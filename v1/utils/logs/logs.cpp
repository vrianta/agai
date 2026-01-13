#include "logs.h"
#include "../../config/config.h"
#include <chrono>
#include <iostream>
#include <cstdarg>

namespace Agai::Utils {

void log(const char *msg) {
  if (Agai::GetConfig().EnableLog == false) {
    return;
  }
  auto now = std::chrono::steady_clock::now().time_since_epoch();
  auto ms = std::chrono::duration_cast<std::chrono::milliseconds>(now).count();
  std::cerr << "[" << ms << " ms] " << msg << std::endl;
}

void logf(const char *fmt, ...) {
  if (!Agai::GetConfig().EnableLog) {
    return;
  }

  auto now = std::chrono::steady_clock::now().time_since_epoch();
  auto ms = std::chrono::duration_cast<std::chrono::milliseconds>(now).count();

  std::cerr << "[" << ms << " ms] ";

  va_list args;
  va_start(args, fmt);
  vfprintf(stderr, fmt, args);
  va_end(args);

  std::cerr << '\n';
}

template <typename... Args> 
void Logln(Args &&...args) {
  if (!Agai::GetConfig().EnableLog)
    return;

  auto now = std::chrono::steady_clock::now().time_since_epoch();
  auto ms = std::chrono::duration_cast<std::chrono::milliseconds>(now).count();

  std::cerr << "[" << ms << " ms] ";

  ((std::cerr << std::forward<Args>(args) << ' '), ...);

  std::cerr << '\n';
}

} // namespace Agai::Utils