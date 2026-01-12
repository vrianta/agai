#ifndef PERSING_CPP
#define PERSING_CPP

#include <string_view>
#include <unordered_map>
#include "agai.cpp"


static Agai::HttpMethod parse_method(std::string_view m);
static void parse_query(std::string_view q, std::unordered_map<std::string_view,std::string_view>& out);
static Agai::HttpRequest parse_request(char *buf, size_t len);


static void parse_query(std::string_view q,
                        std::unordered_map<std::string_view,std::string_view>& out) {
  while (!q.empty()) {
    auto eq = q.find('=');
    auto amp = q.find('&');
    if (eq == std::string_view::npos) break;

    out[q.substr(0, eq)] =
        q.substr(eq + 1, amp - eq - 1);

    if (amp == std::string_view::npos) break;
    q.remove_prefix(amp + 1);
  }
}

static Agai::HttpRequest parse_request(char *buf, size_t len) {
  Agai::HttpRequest req;
  std::string_view data(buf, len);

  // request line
  auto line_end = data.find("\r\n");
  auto line = data.substr(0, line_end);
  data.remove_prefix(line_end + 2);

  auto sp1 = line.find(' ');
  auto sp2 = line.find(' ', sp1 + 1);

  req.method = parse_method(line.substr(0, sp1));
  req.http_version = line.substr(sp2 + 1);

  auto path = line.substr(sp1 + 1, sp2 - sp1 - 1);
  auto qpos = path.find('?');
  if (qpos != std::string_view::npos) {
    req.path = path.substr(0, qpos);
    parse_query(path.substr(qpos + 1), req.query);
  } else {
    req.path = path;
  }

  // headers
  while (true) {
    auto eol = data.find("\r\n");
    if (eol == 0) {
      data.remove_prefix(2);
      break;
    }

    auto h = data.substr(0, eol);
    data.remove_prefix(eol + 2);

    auto colon = h.find(':');
    auto key = h.substr(0, colon);
    auto val = h.substr(colon + 2);
    req.headers[key] = val;

    if (key == "Cookie") {
      while (!val.empty()) {
        auto eq = val.find('=');
        auto sc = val.find(';');
        req.cookies[val.substr(0, eq)] = val.substr(eq + 1, sc - eq - 1);
        if (sc == std::string_view::npos)
          break;
        val.remove_prefix(sc + 2);
      }
    }
  }

  req.body = data;
  return req;
}

static Agai::HttpMethod parse_method(std::string_view m) {
  if (m == "GET")     return Agai::HttpMethod::GET;
  if (m == "POST")    return Agai::HttpMethod::POST;
  if (m == "PUT")     return Agai::HttpMethod::PUT;
  if (m == "DELETE")  return Agai::HttpMethod::DELETE;
  if (m == "PATCH")   return Agai::HttpMethod::PATCH;
  if (m == "OPTIONS") return Agai::HttpMethod::OPTIONS;
  if (m == "HEAD")    return Agai::HttpMethod::HEAD;
  return Agai::HttpMethod::UNKNOWN;
}


#endif