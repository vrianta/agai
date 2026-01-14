#ifndef HEADER_SERVER
#define HEADER_SERVER

#include <string_view>
#include <unordered_map>
#include <vector>

#include "uploads.h"

namespace Agai {
enum HttpMethod : int { GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD, UNKNOWN };
struct HttpRequest {
  HttpMethod method;
  std::string_view path;
  std::string_view http_version;
  std::string_view body;

  std::unordered_map<std::string_view, std::string_view> headers;
  std::unordered_map<std::string_view, std::string_view> query;
  std::unordered_map<std::string_view, std::string_view> cookies;
  std::vector<UploadedFile> files; // Uploaded files
};
} // namespace Agai

#endif