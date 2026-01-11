#ifndef COMPONENTS_CPP
#define COMPONENTS_CPP

#include "agai.h"
#include <functional>
#include <string_view>

namespace {
using Handler = std::function<std::string&(const Agai::HttpRequest &)>;
static std::unordered_map<std::string_view, Handler> get_routes_;
static std::unordered_map<std::string, std::string> templates;
};// namespace

#endif