// routes.h
#pragma once
#include <unordered_map>
#include <string>
#include <functional>
#include "../../components.cpp"

namespace Agai {
using Handler = std::function<Response(const HttpRequest &)>;

bool Get(const std::string& path, Handler handler);
bool Post(const std::string& path, Handler handler);
bool Put(const std::string& path, Handler handler);
bool Delete(const std::string& path, Handler handler);
bool Patch(const std::string& path, Handler handler);

Response RunRequest(const HttpRequest &req);

}
