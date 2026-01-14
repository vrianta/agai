#ifndef HEADER_CONFIG
#define HEADER_CONFIG

#include <string>

namespace Agai {
struct AppSettings {
  bool EnableLog = false;
  std::string StaticFilesDir = "public";
  std::string ViewDirectory = "./src/views";
  int Port = 8080;
};

// will be defined in the main programe
void ConfigSetup(AppSettings &config);
void InitConfig();
const AppSettings &GetConfig();

} // namespace Agai

#endif