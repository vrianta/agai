#ifndef HEADER_UPLOADS
#define HEADER_UPLOADS

#include <string_view>

namespace Agai {

// Structure to represent an uploaded file
struct UploadedFile {
  std::string_view filename;   // Original filename (e.g., "image.png")
  std::string_view mime_type;  // MIME type (e.g., "image/png")
  std::string_view content;    // Raw file data
  std::string_view field_name; // Form field name (e.g., "profileImage")
};

} // namespace Agai

#endif