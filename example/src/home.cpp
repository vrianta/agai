#include "../vendor/agai/include/agai.h"

auto Home = Agai::Get("/", [](const Agai::HttpRequest &req) -> std::string {
  return "Hello from home";
});
