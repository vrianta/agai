#include "../../vendor/agai/include/agai.h"

auto Home = Agai::Get("/", [](const Agai::HttpRequest &req) -> const std::string& {
 return Agai::View("home");
});
