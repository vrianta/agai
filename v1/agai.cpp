#ifndef AGAI_CPP
#define AGAI_CPP

#include <string>
#include <string_view>

#include "agai.h"
#include "components.cpp"
#include "response.h"


// ------------------ API ------------------

bool Agai::Get(const std::string &path, Handler handler) {
  get_routes_[path] = handler;
  return true;
}

Agai::Response Agai::View(const std::string& view) {
    auto it = templates.find(view);
    if (it == templates.end()) {
        return Agai::EmptyResponse; // or throw / return 404 page
    }
    return it->second;
}

const Agai::AppSettings& Agai::GetConfig(){
  return appSettings;
}

#endif