#include "../../response/response.h"
#include <string_view>

namespace Agai {

struct RedirectStatus {
  static constexpr const char *MultipleChoices = "300 Multiple Choices";
  static constexpr const char *MovedPermanently = "301 Moved Permanently";
  static constexpr const char *Found = "302 Found";
  static constexpr const char *SeeOther = "303 See Other";
  static constexpr const char *NotModified = "304 Not Modified";
  static constexpr const char *TemporaryRedirect = "307 Temporary Redirect";
  static constexpr const char *PermanentRedirect = "308 Permanent Redirect";
};

Response Redirect(std::string_view path) {
  Response r("");
  r.SetStatus(RedirectStatus::MovedPermanently);
  r.AddHeader("Location", path);
  return r;
}

} // namespace Agai
