#ifndef HEADER_AGAI_H
#define HEADER_AGAI_H

#include <functional>
#include <string>
#include <string_view>

namespace Agai {

enum class HttpMethod {
  GET,
  POST,
  PUT,
  DELETE_,
  PATCH,
  OPTIONS,
  HEAD,
  UNKNOWN
};
struct HttpRequest {
  HttpMethod method;
  std::string_view path;
  std::string_view http_version;
  std::string_view body;

  std::unordered_map<std::string_view, std::string_view> headers;
  std::unordered_map<std::string_view, std::string_view> query;
  std::unordered_map<std::string_view, std::string_view> cookies;
};
// React-like API
const std::string Get(const std::string &, std::function<std::string&(const HttpRequest &)>);

// base folder will be Views where the executable present and the view strign will be folder.file_name or file_name directly
// 
std::string View(std::string_view view);
}; // namespace Agai
#endif
