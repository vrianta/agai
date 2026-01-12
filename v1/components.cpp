#ifndef COMPONENTS_CPP
#define COMPONENTS_CPP

#include "response.h"
#include "agai.h"

#include <functional>
#include <string>

namespace {
using Handler = std::function<Agai::Response(const Agai::HttpRequest &)>;
static std::unordered_map<std::string, Handler> get_routes_ = {
  {
    "/404",
    [](const Agai::HttpRequest&) {
      return Agai::View("home");
    }
  }
};

static std::unordered_map<std::string, Agai::Response> templates;

auto appSettings = Agai::AppSettings{};

};// namespace

#endif