#include "utils.h"

// ------------------ utils ------------------

static std::vector<std::string> split(const std::string &s, char delim) {
  std::vector<std::string> out;
  std::stringstream ss(s);
  std::string item;
  while (std::getline(ss, item, delim))
    out.push_back(item);
  return out;
}

static std::string trim(std::string s) {
  while (!s.empty() && isspace(s.front()))
    s.erase(0, 1);
  while (!s.empty() && isspace(s.back()))
    s.pop_back();
  return s;
}