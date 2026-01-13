#ifndef ROUTING_CPP
#define ROUTING_CPP

// routes.cpp
#include "routing.h"
#include "../../404.cpp"
#include "../../agai.h"
#include "../../response/response.h"
#include "../redirects/redirects.cpp"

namespace Agai {

static constexpr const char *NOT_FOUND_PATH = "/404";
static constexpr const char *INTERNAL_ERROR_PATH = "/500";
static const Response NotFound_Response = Response(_404);
static const Response Note_Found_redirect = Redirect("/404");
static Handler _404_Handler = [](const Agai::HttpRequest &) {
  return NotFound_Response;
};

static std::unordered_map<std::string, Handler> get_routes = {
    {NOT_FOUND_PATH, _404_Handler}};
static std::unordered_map<std::string, Handler> post_routes{};
static std::unordered_map<std::string, Handler> put_routes{};
static std::unordered_map<std::string, Handler> delete_routes{};
static std::unordered_map<std::string, Handler> patch_routes{};

bool Get(const std::string &path, Handler handler) {
  get_routes[path] = std::move(handler);
  return true;
}

bool Post(const std::string &path, Handler handler) {
  post_routes[path] = std::move(handler);
  return true;
}

bool Put(const std::string &path, Handler handler) {
  put_routes[path] = std::move(handler);
  return true;
}

bool Delete(const std::string &path, Handler handler) {
  delete_routes[path] = std::move(handler);
  return true;
}

bool Patch(const std::string &path, Handler handler) {
  patch_routes[path] = std::move(handler);
  return true;
}

Response RunRequest(const HttpRequest &req) {
  switch (req.method) {
  case HttpMethod::GET: {
    auto it = get_routes.find(std::string(req.path));
    if (it != get_routes.end()) {
      return it->second(req);
    } else {
      return Note_Found_redirect;
    }
    break;
  }
  case HttpMethod::POST: {
    auto it = post_routes.find(std::string(req.path));
    if (it != get_routes.end()) {
      return it->second(req);
    } else {
      return Note_Found_redirect;
    }
    break;
  }
  case HttpMethod::PUT: {
    auto it = put_routes.find(std::string(req.path));
    if (it != get_routes.end()) {
      return it->second(req);
    } else {
      return Note_Found_redirect;
    }
    break;
  }
  case HttpMethod::DELETE: {
    auto it = delete_routes.find(std::string(req.path));
    if (it != get_routes.end()) {
      return it->second(req);
    } else {
      return Note_Found_redirect;
    }
    break;
  }
  case HttpMethod::PATCH: {
    auto it = patch_routes.find(std::string(req.path));
    if (it != get_routes.end()) {
      return it->second(req);
    } else {
      return Note_Found_redirect;
    }
    break;
  }
  default:
    return Note_Found_redirect;
  }
  return EmptyResponse;
}

} // namespace Agai

#endif