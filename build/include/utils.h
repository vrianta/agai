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
#include "../agai.h"

namespace Agai::Utils {

static std::vector<std::string> split(const std::string &s, char delim);
static std::string trim(std::string s);

// File operations
bool saveFile(const UploadedFile &file, const std::string &directory = "");
bool saveFileToPath(const UploadedFile &file, const std::string &fullPath);
bool saveFiles(const std::vector<UploadedFile> &files, const std::string &directory = "");

} // namespace Utils

#endif
