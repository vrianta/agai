#include "agai.h"
#include "utils.h"

#include <chrono>
#include <iostream>


// ------------------ utils ------------------

static std::vector<std::string> Agai::Utils::split(const std::string &s, char delim) {
  std::vector<std::string> out;
  std::stringstream ss(s);
  std::string item;
  while (std::getline(ss, item, delim))
    out.push_back(item);
  return out;
}

static std::string Agai::Utils::trim(std::string s) {
  while (!s.empty() && isspace(s.front()))
    s.erase(0, 1);
  while (!s.empty() && isspace(s.back()))
    s.pop_back();
  return s;
}


void Agai::Utils::log(const char* msg) {
  if (Agai::GetConfig().EnableLog == false) {
    return;
  }
  auto now = std::chrono::steady_clock::now().time_since_epoch();
  auto ms = std::chrono::duration_cast<std::chrono::milliseconds>(now).count();
  std::cerr << "[" << ms << " ms] " << msg << std::endl;
}

void Agai::Utils::logf(const char* fmt, ...)
{
    if (!Agai::GetConfig().EnableLog) {
        return;
    }

    auto now = std::chrono::steady_clock::now().time_since_epoch();
    auto ms = std::chrono::duration_cast<std::chrono::milliseconds>(now).count();

    std::cerr << "[" << ms << " ms] ";

    va_list args;
    va_start(args, fmt);
    vfprintf(stderr, fmt, args);
    va_end(args);

    std::cerr << '\n';
}

void Logln(const char* msg)
{
    if (!Agai::GetConfig().EnableLog) {
        return;
    }

    auto now = std::chrono::steady_clock::now().time_since_epoch();
    auto ms = std::chrono::duration_cast<std::chrono::milliseconds>(now).count();

    std::cerr << "[" << ms << " ms] " << msg << '\n';
}

template<typename... Args>
void Logln(Args&&... args)
{
    if (!Agai::GetConfig().EnableLog) return;

    auto now = std::chrono::steady_clock::now().time_since_epoch();
    auto ms = std::chrono::duration_cast<std::chrono::milliseconds>(now).count();

    std::cerr << "[" << ms << " ms] ";

    ((std::cerr << std::forward<Args>(args) << ' '), ...);

    std::cerr << '\n';
}
