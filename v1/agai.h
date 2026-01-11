#ifndef HEADER_AGAI_H
#define HEADER_AGAI_H

#include <functional>
#include <string>
#include <string_view>

namespace Agai {

struct AppSettings {
  bool EnableLog = false;
  int Port = 8080;
};

// will be defined in the main programe
void ConfigSetup(AppSettings &config);
const AppSettings &GetConfig();

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
bool Get(const std::string &,
         std::function<const std::string &(const HttpRequest &)>);

// it will check get the template with view index and return it
const std::string View(std::string_view view);

// Function to setup the
}; // namespace Agai
#endif
