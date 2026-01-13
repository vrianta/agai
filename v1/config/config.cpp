#include "config.h"

namespace Agai {


static auto appSettings = Agai::AppSettings{};

const Agai::AppSettings &GetConfig() { return appSettings; }

void InitConfig() {
    ConfigSetup(appSettings);
}

} // namespace Agai