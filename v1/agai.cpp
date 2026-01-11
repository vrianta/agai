#ifndef AGAI_CPP
#define AGAI_CPP

#include <string>
#include <string_view>

#include "agai.h"
#include "components.cpp"


// ------------------ API ------------------

bool Agai::Get(const std::string &path, Handler handler) {
  get_routes_[path] = handler;
  return true;
}

const std::string Agai::View(std::string_view view) {
    auto it = templates.find(std::string(view));
    if (it == templates.end()) {
        return {}; // or throw / return 404 page
    }
    return it->second;
}

const Agai::AppSettings& Agai::GetConfig(){
  return appSettings;
}

#endif