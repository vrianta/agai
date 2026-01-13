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


#include "config/config.h"
#include "response/response.h"
#include "server/server.cpp"

void RegisterTemplates();

int main() {
  Agai::InitConfig();
  RegisterTemplates();
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
  Agai::Utils::logf("RegisterTemplates: start");

  const fs::path root = "Views";
  Agai::Utils::logf("RegisterTemplates: root path = %s", root.string().c_str());

  templates.clear();
  Agai::Utils::logf("RegisterTemplates: templates map cleared");

  for (const auto &entry : fs::recursive_directory_iterator(root)) {
    if (!entry.is_regular_file()) {
      Agai::Utils::logf(
        "RegisterTemplates: skipped non-regular entry: %s",
        entry.path().string().c_str()
      );
      continue;
    }

    const fs::path &path = entry.path();
    Agai::Utils::logf(
      "RegisterTemplates: processing file: %s",
      path.string().c_str()
    );

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

    Agai::Utils::logf(
      "RegisterTemplates: generated key = %s",
      key.c_str()
    );

    // read file
    std::ifstream file(path, std::ios::binary);
    if (!file) {
      Agai::Utils::logf(
        "RegisterTemplates: ERROR failed to open file: %s",
        path.string().c_str()
      );
      continue;
    }

    std::ostringstream ss;
    ss << file.rdbuf();
    const std::string content = ss.str();

    Agai::Utils::logf(
      "RegisterTemplates: read %zu bytes from %s",
      content.size(),
      path.string().c_str()
    );

    templates.emplace(key, Agai::Response(content.c_str()));

    Agai::Utils::logf(
      "RegisterTemplates: template registered: key=%s",
      key.c_str()
    );
  }

  Agai::Utils::logf(
    "RegisterTemplates: completed, total templates=%zu",
    templates.size()
  );
}

