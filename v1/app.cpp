// file: http_server.cpp

#include <arpa/inet.h>
#include <cstring>
#include <string>
#include <sys/socket.h>
#include <unistd.h>
#include <unordered_map>

#include "config/config.h"
#include "response/response.h"
#include "server/server.cpp"
#include "components.cpp"

namespace Agai {

using views_map = std::map<std::string, std::vector<unsigned char>>;

extern views_map register_embedded_views();

} // namespace Agai

void RegisterTemplates(const Agai::views_map& views);

int main() {
  Agai::InitConfig();
  RegisterTemplates(Agai::register_embedded_views());
  serve("0.0.0.0", 8080);
}

/*
 * Look int to the folder and recursively check all the folder and files
 * Store the template contentns in a map with the index which would look like
 * this if the template is in view directory then it will be file name without
 * extension if the file is any folder then index will be
 * folder_name.folder_name.file_name
 */

void RegisterTemplates(const Agai::views_map& views) {
  templates.clear();

  Agai::Utils::logf("[Templates] registering %zu templates", views.size());

  for (const auto& view : views) {
    const auto& name = view.first;
    const auto& data = view.second;

    Agai::Utils::logf(
      "  - %s (%zu bytes)",
      name.c_str(),
      data.size()
    );

    templates.emplace(
      name,
      Agai::Response(std::string(
        reinterpret_cast<const char*>(data.data()),
        data.size()
      ))
    );
  }

  Agai::Utils::logf("[Templates] registered templates:");
  for (const auto& t : templates) {
    Agai::Utils::logf("  * %s", t.first.c_str());
  }
}

