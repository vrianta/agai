
#include "agai.h"
#include "utils.h"

#include <arpa/inet.h>
#include <iostream>
#include <sys/socket.h>
#include <unistd.h>

#include <cstring>
#include <string>
#include <unordered_map>

#include "components.cpp"
#include "persing.cpp"

std::string not_found = "Not Found";

static Agai::HttpRequest parse_request(char *buf, size_t len);

static int serve(const char *host, int port) {
  Agai::Utils::log("serve() entered");

  Agai::Utils::log("creating socket");
  int s = socket(AF_INET, SOCK_STREAM, 0);

  Agai::Utils::log("setting socket options");
  int opt = 1;
  setsockopt(s, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));

  Agai::Utils::log("preparing sockaddr");
  sockaddr_in addr{};
  addr.sin_family = AF_INET;
  addr.sin_port = htons(port);
  inet_pton(AF_INET, host, &addr.sin_addr);

  Agai::Utils::log("binding socket");
  if (bind(s, (sockaddr *)&addr, sizeof(addr)) < 0) {
    perror("bind failed");
    close(s);
    std::exit(EXIT_FAILURE);
  }

  Agai::Utils::log("listening");
  if (listen(s, 10) < 0) {
    perror("listen failed");
    close(s);
    std::exit(EXIT_FAILURE);
  }

  Agai::Utils::log("server ready");
  std::cout << "Server started on " << host << ":" << port << std::endl;
  std::string status = "200 OK";

  while (true) {
    int c = accept(s, nullptr, nullptr);
    char buffer[8192];
    ssize_t n = read(c, buffer, sizeof(buffer));
    auto req = parse_request(buffer, n);
    std::string body;
    switch (req.method) {
    case Agai::HttpMethod::GET: {
      auto it = get_routes_.find(req.path);
      if (it != get_routes_.end())
        body = it->second(req);
      else {
        status = "404 Not Found";
        body = "Route Not Found";
      }
      break;
    }
    case Agai::HttpMethod::POST:
      break;
    case Agai::HttpMethod::PUT:
      break;
    case Agai::HttpMethod::DELETE_:
      break;
    case Agai::HttpMethod::PATCH:
      break;
    case Agai::HttpMethod::OPTIONS:
      break;
    case Agai::HttpMethod::HEAD:
      break;
    case Agai::HttpMethod::UNKNOWN:
      break;
    default:
      body = "Method Not Allowed";
      break;
    }
    std::string res;
    res.reserve(128 + body.size());

    res.append("HTTP/1.1 ");
    res.append(status);
    res.append(
        "\r\nContent-Type: text/html; charset=utf-8\r\nContent-Length: ");
    res.append(std::to_string(body.size()));
    res.append("\r\n\r\n");
    res.append(body);

    write(c, res.c_str(), res.size());
    close(c);
  }
}
