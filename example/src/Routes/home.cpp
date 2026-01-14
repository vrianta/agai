#include "../../vendor/agai/include/v1/response/response.h"
#include "../../vendor/agai/include/v1/agai.h"
#include "../../vendor/agai/include/v1/server/routing/routing.h"

auto Home = Agai::Get("/", [](const Agai::HttpRequest &req) -> Agai::Response {
  return Agai::View("themes.atlas-portfolio.home");
});
