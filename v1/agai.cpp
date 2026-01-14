#ifndef AGAI_CPP
#define AGAI_CPP

#include <string>
#include <map>

#include "agai.h"
#include "components.cpp"
#include "response/response.h"

namespace Agai {

// extern std::map<std::string, std::vector<unsigned char>> register_embedded_views();

Response View(const std::string &view) {
  auto it = templates.find(view);
  if (it == templates.end()) {
    return Agai::EmptyResponse; // or throw / return 404 page
  }
  return it->second;
}
} // namespace Agai

#endif