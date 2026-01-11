// file: http_server.cpp

#include <arpa/inet.h>
#include <sys/socket.h>
#include <unistd.h>
#include <cstring>
#include <string>
#include <unordered_map>
#include <filesystem>
#include <fstream>
#include <sstream>


#include "agai.h"
#include "utils.cpp"
#include "server.cpp"

void RegisterTemplates();

int main() {
  RegisterTemplates();
  Agai::ConfigSetup(appSettings);
  serve("0.0.0.0", 8080);
}

/*
 * Look int to the folder and recursively check all the folder and files
 * Store the template contentns in a map with the index which would look like
 * this if the template is in view directory then it will be file name without
 * extension if the file is any folder then index will be
 * folder_name.folder_name.file_name
 */

namespace fs = std::filesystem;

void RegisterTemplates() {
  const fs::path root = "Views";

  templates.clear();

  for (const auto &entry : fs::recursive_directory_iterator(root)) {
    if (!entry.is_regular_file())
      continue;

    const fs::path &path = entry.path();

    // build key
    fs::path rel = fs::relative(path, root);
    std::string key;

    for (auto it = rel.begin(); it != rel.end(); ++it) {
      if (it->has_extension()) {
        key += it->stem().string();
      } else {
        key += it->string();
      }
      if (std::next(it) != rel.end())
        key += ".";
    }

    // read file
    std::ifstream file(path, std::ios::binary);
    std::ostringstream ss;
    ss << file.rdbuf();

    templates.emplace(std::move(key), ss.str());
  }
}
