#ifndef AGAI_CPP
#define AGAI_CPP

#include "agai.h"
#include "components.cpp"
#include <string>
#include <string_view>

// ------------------ API ------------------

const std::string Agai::Get(const std::string &path, Handler handler) {
  get_routes_[path] = handler;
  return std::string("");
}

std::string Agai::View(std::string_view view) {

}

#endif