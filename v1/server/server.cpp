
#include <arpa/inet.h>
#include <cstdlib>
#include <iostream>
#include <sys/socket.h>
#include <unistd.h>

#include <cstring>
#include <string>
#include <unordered_map>

#include "../components.cpp"
#include "../persing.cpp"
#include "../agai.h"
#include "../utils.h"
#include "redirects.cpp"

std::string not_found = "Not Found";
static const char *HttpMethodString[] = {"Get",   "Post",    "Put",  "Delete",
                                         "Patch", "Options", "Head", "Unknown"};

static Agai::HttpRequest parse_request(char *buf, size_t len);

static int serve(const char *host, int port) {
  Agai::Utils::logf("[server] starting (host=%s port=%d)", host, port);

  int s = socket(AF_INET, SOCK_STREAM, 0);
  if (s < 0) {
    perror("socket");
    return -1;
  }
  Agai::Utils::logf("[server] socket created fd=%d", s);

  int opt = 1;
  setsockopt(s, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));
  Agai::Utils::log("[server] socket options set (SO_REUSEADDR)");

  sockaddr_in addr{};
  addr.sin_family = AF_INET;
  addr.sin_port = htons(port);
  inet_pton(AF_INET, host, &addr.sin_addr);

  if (bind(s, (sockaddr *)&addr, sizeof(addr)) < 0) {
    perror("bind");
    close(s);
    return -1;
  }
  Agai::Utils::log("[server] bind successful");

  if (listen(s, 10) < 0) {
    perror("listen");
    close(s);
    return -1;
  }
  Agai::Utils::log("[server] listening");

  std::cout << "Server started on " << host << ":" << port << std::endl;

  while (true) {
    Agai::Utils::log("[server] waiting for connection");
    int c = accept(s, nullptr, nullptr);
    if (c < 0) {
      perror("accept");
      continue;
    }
    Agai::Utils::logf("[conn] accepted fd=%d", c);

    // Dynamic buffer handling with growth
    std::string buffer;
    const size_t INITIAL_BUFFER_SIZE = 4096;
    const size_t MAX_BUFFER_SIZE = 10 * 1024 * 1024; // 10MB limit
    buffer.reserve(INITIAL_BUFFER_SIZE);

    ssize_t total_read = 0;
    bool headers_complete = false;
    size_t header_end = 0;
    bool request_complete = false;

    // Read until we have complete headers and full body (if any)
    while (total_read < MAX_BUFFER_SIZE && !request_complete) {
      char temp_buffer[8192];
      ssize_t n = read(c, temp_buffer, sizeof(temp_buffer));
      
      if (n < 0) {
        perror("read");
        close(c);
        continue;
      }
      if (n == 0) {
        // Connection closed by client
        if (total_read == 0) {
          Agai::Utils::logf("[conn] read failed fd=%d bytes=0", c);
          close(c);
          request_complete = true; // Force exit
          continue;
        }
        break;
      }

      buffer.append(temp_buffer, n);
      total_read += n;

      // Find headers end (\r\n\r\n)
      if (!headers_complete) {
        header_end = buffer.find("\r\n\r\n");
        if (header_end != std::string::npos) {
          headers_complete = true;
          header_end += 4; // Include the \r\n\r\n

          // Check Content-Length to know if we need to read more
          size_t content_length_pos = buffer.find("Content-Length:");
          if (content_length_pos != std::string::npos) {
            content_length_pos += 15; // strlen("Content-Length:")
            size_t eol = buffer.find("\r\n", content_length_pos);
            std::string_view length_str = std::string_view(buffer).substr(
                content_length_pos, eol - content_length_pos);
            // Trim whitespace
            while (!length_str.empty() && length_str[0] == ' ')
              length_str.remove_prefix(1);
            
            size_t content_length = std::stoul(std::string(length_str));
            size_t body_size = total_read - header_end;
            
            if (body_size >= content_length) {
              request_complete = true;
            }
            // Continue reading to get full body
          } else {
            // No Content-Length header, assume request is complete after headers
            request_complete = true;
          }
        }
      }

      // If buffer is getting too large, break to avoid DoS
      if (buffer.size() > MAX_BUFFER_SIZE) {
        Agai::Utils::logf("[conn] buffer size exceeded fd=%d", c);
        close(c);
        request_complete = true;
        continue;
      }
    }

    if (total_read == 0 || !request_complete) {
      Agai::Utils::logf("[conn] read failed fd=%d bytes=%zd", c, total_read);
      close(c);
      continue;
    }

    Agai::Utils::logf("[conn] received %zd bytes", total_read);

    auto req = parse_request((char*)buffer.c_str(), buffer.size());
    Agai::Utils::logf("[http] %s %s",
                      HttpMethodString[int(req.method)], // assume helper
                      req.path);

    Agai::Response res;
    switch (req.method) {
    case Agai::HttpMethod::GET: {
      auto it = get_routes_.find(req.path.data());
      if (it != get_routes_.end()) {
        Agai::Utils::logf("[router] matched GET %s", req.path);
        res = it->second(req);
      } else {
        Agai::Utils::logf("[router] no route for GET %s", req.path);
        res = Agai::Redirect("/404");
      }
      break;
    }

    default:
      Agai::Utils::logf("[http] %s %s",
                        HttpMethodString[int(req.method)], // assume helper
                        req.path);
      res = Agai::Redirect("/404");
      break;
    }
    auto content = res.Serialize();
    ssize_t written = write(c, content.c_str(), content.size());

    close(c);
    Agai::Utils::logf("[conn] closed fd=%d", c);
  }
}
