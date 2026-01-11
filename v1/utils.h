#include <arpa/inet.h>
#include <sys/socket.h>
#include <unistd.h>
#include <vector>

#include <cstring>
#include <sstream>
#include <string>
#include <vector>

namespace Agai {
namespace Utils {

static std::vector<std::string> split(const std::string &s, char delim);
static std::string trim(std::string s);
} // namespace Utils

} // namespace Agai
