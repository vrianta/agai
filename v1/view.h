namespace Agai {
class View {
private:
  enum class Types {
    html,
    chtml,
  };
  const char *type_as_string;
  char buffer[8192];
  Types type; // store the template type as

public:
  // pass the type as string
  const char *GetType();
};
} // namespace Agai