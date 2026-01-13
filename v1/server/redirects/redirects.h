#ifndef REDIRECTS_HEADER
#define REDIRECTS_HEADER

#include "../../response/response.h"

namespace Agai {
Response Redirect(std::string_view path);
}
#endif