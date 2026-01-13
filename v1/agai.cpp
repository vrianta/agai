#ifndef AGAI_CPP
#define AGAI_CPP

#include <string>

#include "agai.h"
#include "components.cpp"
#include "response/response.h"



Agai::Response Agai::View(const std::string& view) {
    auto it = templates.find(view);
    if (it == templates.end()) {
        return Agai::EmptyResponse; // or throw / return 404 page
    }
    return it->second;
}

#endif