#ifndef Response_HEADER
#define Response_HEADER

#include <cstdint>
#include <string>
#include <string_view>

namespace Agai {

class Response {
private:
  enum Types : int {
    html,
    json,
  };
  static constexpr const char *ContentTypes[] = {
      [Types::html] = "\r\ntext/html; charset=utf-8\r\nContent-Length: ",
      [Types::json] = "\r\napplication/json\r\nContent-Length: "};

  std::string body;
  std::string status = "200 OK";
  int64_t Length;
  Types type = Types::html; // store the Response type as
  std::string headers; // owns all extra headers

public:
  Response(std::string content);
  Response() = default;
  ~Response() = default;
  // pass the type as string
  const char *GetType() const;
  void AsJson();
  const char *GetContent() const;
  void AddHeader(std::string_view key, std::string_view value);

  void SetStatus(std::string s);
  std::string Serialize() const;
};

// header
inline Agai::Response EmptyResponse("");

}; // namespace Agai

#endif