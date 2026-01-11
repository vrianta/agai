#ifndef UTILS_HEADER
#define UTILS_HEADER
#include <arpa/inet.h>
#include <sys/socket.h>
#include <unistd.h>
#include <vector>

#include <cstring>
#include <sstream>
#include <string>
#include <vector>

namespace Agai::Utils {

static std::vector<std::string> split(const std::string &s, char delim);
static std::string trim(std::string s);
void log(const char *msg);
void logf(const char *msg, ...);
template<typename... Args>
void Logln(Args&&... args);

} // namespace Utils

#endif
