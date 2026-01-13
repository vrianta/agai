#ifndef HEADER_AGAI_H
#define HEADER_AGAI_H

#include "response/response.h"
#include <functional>
#include <string>
#include <string_view>
#include <vector>

namespace Agai {

// Structure to represent an uploaded file
struct UploadedFile {
  std::string_view filename;        // Original filename (e.g., "image.png")
  std::string_view mime_type;       // MIME type (e.g., "image/png")
  std::string_view content;         // Raw file data
  std::string_view field_name;      // Form field name (e.g., "profileImage")
};

enum HttpMethod : int {
  GET,
  POST,
  PUT,
  DELETE,
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
  std::vector<UploadedFile> files;                                       // Uploaded files
};

// it will check get the template with view index and return it
Agai::Response View(const std::string& view);

// Function to setup the
}; // namespace Agai
#endif
