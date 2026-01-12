#include "../../vendor/agai/include/agai.h"

auto Home = Agai::Get("/", [](const Agai::HttpRequest &req) -> Agai::Response {
  return Agai::View("home");
});
