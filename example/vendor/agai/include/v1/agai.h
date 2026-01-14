#ifndef HEADER_AGAI_H
#define HEADER_AGAI_H


#include <functional>
#include <string>
#include <string_view>
#include <vector>

#include "response/response.h"
#include "config/config.h"
#include "server/routing/routing.h"

namespace Agai {

// it will check get the template with view index and return it
Agai::Response View(const std::string& view);

// Function to setup the
}; // namespace Agai
#endif
