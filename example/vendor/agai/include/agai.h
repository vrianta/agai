#ifndef HEADER_AGAI_H
#define HEADER_AGAI_H

#include <functional>
#include <string>

namespace Agai {
struct HttpRequest {
  std::string_view method;
  std::string_view path;
  std::string_view http_version;
  std::string_view body;

  std::unordered_map<std::string_view, std::string_view> headers;
  std::unordered_map<std::string_view, std::string_view> query;
  std::unordered_map<std::string_view, std::string_view> cookies;
};
// React-like API
bool Get(const std::string &, std::function<std::string(const HttpRequest &)>);

std::string& View();
}; // namespace Agai
#endif
