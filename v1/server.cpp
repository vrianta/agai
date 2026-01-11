
#include "agai.h"

#include <arpa/inet.h>
#include <string_view>
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
  int s = socket(AF_INET, SOCK_STREAM, 0);
  int opt = 1;
  setsockopt(s, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));

  sockaddr_in addr{};
  addr.sin_family = AF_INET;
  addr.sin_port = htons(port);
  inet_pton(AF_INET, host, &addr.sin_addr);

  if (bind(s, (sockaddr *)&addr, sizeof(addr)) < 0) {
    perror("bind failed");
    close(s);
    std::exit(EXIT_FAILURE);
  }

  if (listen(s, 10) < 0) {
    perror("listen failed");
    close(s);
    std::exit(EXIT_FAILURE);
  }

  while (true) {
    int c = accept(s, nullptr, nullptr);
    char buffer[8192];
    ssize_t n = read(c, buffer, sizeof(buffer));

    auto req = parse_request(buffer, n);

    std::string body;

    switch (req.method) {

    case Agai::HttpMethod::GET:
      if (get_routes_.count(std::string(req.path)))
        body = get_routes_[req.path](req);
      else
        body = "Route Not Found";
      break;
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

    std::string res = "HTTP/1.1 200 OK\r\n"
                      "Content-Length: " +
                      std::to_string(body.size()) + "\r\n\r\n" + body;

    write(c, res.c_str(), res.size());
    close(c);
  }
}
