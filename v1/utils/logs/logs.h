#ifndef HEADER_LOGS
#define HEADER_LOGS

namespace Agai::Utils {

void log(const char *msg);
void logf(const char *msg, ...);
template <typename... Args> void Logln(Args &&...args);

} // namespace Agai::Utils

#endif