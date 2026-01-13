
#include "response.h"
#include <cstdlib>
#include <cstring>

#include "../utils/logs/logs.h"

namespace Agai {

const char *Response::GetType() const {
  return this->ContentTypes[this->type];
}

void Response::AsJson() { this->type = Types::json; }

const char *Agai::Response::GetContent() const { return this->body.data(); }

Response::Response(std::string content) {
  if (content.size() > 8192) {
    Agai::Utils::log("Response content exceeds max size (8192 bytes)");
    throw std::length_error("Response too large");
  }
  body = std::move(content);
}

void Response::SetStatus(std::string s) { this->status = s; }

std::string Response::Serialize() const {
  std::string res;
  res.reserve(128 + body.size());

  res.append("HTTP/1.1 ");
  res.append(status);
  res.append("\r\nContent-Type: ");
  res.append(ContentTypes[type]);
  res.append(std::to_string(body.size()));
  res.append("\r\n");
  res.append(headers);
  res.append("\r\n");
  res.append(body);

  return res;
}

void Response::AddHeader(std::string_view key, std::string_view value) {
  headers.append(key);
  headers.append(": ");
  headers.append(value);
  headers.append("\r\n");
}
} // namespace Agai
