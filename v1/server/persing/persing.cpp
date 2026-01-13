#ifndef PERSING_CPP
#define PERSING_CPP

#include "../../agai.h"
#include "../../agai.cpp"
#include <string_view>
#include <unordered_map>
#include <cctype>

static Agai::HttpMethod parse_method(std::string_view m);
static void parse_query(
            std::string_view q,
            std::unordered_map<std::string_view, 
            std::string_view> &out
          );
static void parse_form_data(
            std::string_view body,
            std::unordered_map<std::string_view, 
            std::string_view> &out
          );
static void parse_multipart_form(
            std::string_view body,
            std::string_view boundary,
            std::unordered_map<std::string_view, std::string_view> &form_fields,
            std::vector<Agai::UploadedFile> &files
          );
static Agai::HttpRequest parse_request(char *buf, size_t len);

// Trim leading whitespace from string_view
static inline std::string_view trim_left(std::string_view s) {
  while (!s.empty() && std::isspace(s[0])) {
    s.remove_prefix(1);
  }
  return s;
}

// Trim trailing whitespace from string_view
static inline std::string_view trim_right(std::string_view s) {
  while (!s.empty() && std::isspace(s[s.length() - 1])) {
    s.remove_suffix(1);
  }
  return s;
}

// Case-insensitive string comparison
static inline bool iequals(std::string_view a, std::string_view b) {
  if (a.size() != b.size()) return false;
  for (size_t i = 0; i < a.size(); i++) {
    if (std::tolower(a[i]) != std::tolower(b[i])) return false;
  }
  return true;
}

static void parse_query(
            std::string_view q,
            std::unordered_map<std::string_view, 
            std::string_view> &out
          ) {
  while (!q.empty()) {
    auto eq = q.find('=');
    auto amp = q.find('&');
    
    // Handle case with no '=' (skip malformed pairs)
    if (eq == std::string_view::npos) {
      break;
    }
    
    // Handle case with no '&' (last parameter)
    size_t value_len = (amp == std::string_view::npos) ? 
                       (q.size() - eq - 1) : (amp - eq - 1);
    
    std::string_view key = q.substr(0, eq);
    std::string_view value = q.substr(eq + 1, value_len);
    
    // Skip empty keys
    if (!key.empty()) {
      out[key] = value;
    }
    
    if (amp == std::string_view::npos) {
      break;
    }
    q.remove_prefix(amp + 1);
  }
}

// Parse form data from POST body (application/x-www-form-urlencoded)
// Format: key1=value1&key2=value2&key3=value3
static void parse_form_data(
            std::string_view body,
            std::unordered_map<std::string_view, 
            std::string_view> &out
          ) {
  // Form data uses the same format as query strings
  parse_query(body, out);
}

// Parse multipart form data (application/multipart/form-data)
// Used for file uploads and mixed text/file submissions
static void parse_multipart_form(
            std::string_view body,
            std::string_view boundary,
            std::unordered_map<std::string_view, std::string_view> &form_fields,
            std::vector<Agai::UploadedFile> &files
          ) {
  std::string boundary_marker = "--" + std::string(boundary);
  size_t pos = 0;
  
  // Skip initial boundary
  pos = body.find(boundary_marker, pos);
  if (pos == std::string_view::npos) return;
  pos += boundary_marker.length();
  
  while (pos < body.size()) {
    // Skip CRLF/LF after boundary
    if (pos < body.size() - 1 && body[pos] == '\r' && body[pos+1] == '\n') {
      pos += 2;
    } else if (pos < body.size() && body[pos] == '\n') {
      pos += 1;
    }
    
    // Check for end boundary
    if (pos + 1 < body.size() && body[pos] == '-' && body[pos+1] == '-') {
      break;
    }
    
    // Find end of headers
    size_t header_end = body.find("\r\n\r\n", pos);
    size_t header_sep = 4;
    if (header_end == std::string_view::npos) {
      header_end = body.find("\n\n", pos);
      header_sep = 2;
      if (header_end == std::string_view::npos) break;
    }
    
    auto headers = body.substr(pos, header_end - pos);
    pos = header_end + header_sep;
    
    // Parse headers
    std::string_view field_name;
    std::string_view filename;
    std::string_view mime_type = "text/plain";
    
    auto name_pos = headers.find("name=\"");
    if (name_pos != std::string_view::npos) {
      name_pos += 6;
      auto name_end = headers.find("\"", name_pos);
      if (name_end != std::string_view::npos) {
        field_name = headers.substr(name_pos, name_end - name_pos);
      }
    }
    
    auto file_pos = headers.find("filename=\"");
    if (file_pos != std::string_view::npos) {
      file_pos += 10;
      auto file_end = headers.find("\"", file_pos);
      if (file_end != std::string_view::npos) {
        filename = headers.substr(file_pos, file_end - file_pos);
      }
    }
    
    auto mime_pos = headers.find("Content-Type: ");
    if (mime_pos != std::string_view::npos) {
      mime_pos += 14;
      auto mime_end = headers.find("\n", mime_pos);
      if (mime_end != std::string_view::npos) {
        mime_type = trim_right(headers.substr(mime_pos, mime_end - mime_pos));
      }
    }
    
    // Find content boundary
    size_t content_end = body.find("\r\n" + boundary_marker, pos);
    if (content_end == std::string_view::npos) {
      content_end = body.find("\n" + boundary_marker, pos);
      if (content_end == std::string_view::npos) break;
    }
    
    auto content = body.substr(pos, content_end - pos);
    
    if (!filename.empty()) {
      Agai::UploadedFile file;
      file.field_name = field_name;
      file.filename = filename;
      file.mime_type = mime_type;
      file.content = content;
      files.push_back(file);
    } else {
      form_fields[field_name] = content;
    }
    
    pos = content_end;
    if (body[pos] == '\r') pos += 2;
    else pos += 1;
    
    pos += boundary_marker.length();
  }
}

static Agai::HttpRequest parse_request(char *buf, size_t len) {
  Agai::HttpRequest req;
  std::string_view data(buf, len);

  // Parse request line
  auto line_end = data.find("\r\n");
  if (line_end == std::string_view::npos) {
    // Malformed request, no CRLF found
    return req;
  }

  auto line = data.substr(0, line_end);
  
  // Validate we have enough data
  if (line_end + 2 > data.size()) {
    return req;
  }
  data.remove_prefix(line_end + 2);

  // Parse method, path, version
  auto sp1 = line.find(' ');
  auto sp2 = line.find(' ', sp1 + 1);

  if (sp1 == std::string_view::npos || sp2 == std::string_view::npos) {
    // Malformed request line
    return req;
  }

  req.method = parse_method(line.substr(0, sp1));
  req.http_version = trim_right(line.substr(sp2 + 1));

  auto path = line.substr(sp1 + 1, sp2 - sp1 - 1);
  auto qpos = path.find('?');
  
  if (qpos != std::string_view::npos) {
    req.path = path.substr(0, qpos);
    parse_query(path.substr(qpos + 1), req.query);
  } else {
    req.path = path;
  }

  // Parse headers
  while (true) {
    auto eol = data.find("\r\n");
    
    if (eol == std::string_view::npos) {
      // End of headers reached abruptly (body without CRLF marker)
      req.body = data;
      break;
    }
    
    if (eol == 0) {
      // Empty line marks end of headers
      if (data.size() >= 2) {
        data.remove_prefix(2);
      }
      req.body = data;
      break;
    }

    auto h = data.substr(0, eol);
    auto colon = h.find(':');
    
    if (colon != std::string_view::npos && colon + 1 < h.size()) {
      auto key = h.substr(0, colon);
      auto val = trim_left(h.substr(colon + 1));
      
      // Store header (using original case for compatibility)
      req.headers[key] = val;

      // Parse cookies if this is Cookie header (case-insensitive check)
      if (iequals(key, "Cookie")) {
        while (!val.empty()) {
          auto eq = val.find('=');
          auto sc = val.find(';');
          
          if (eq != std::string_view::npos) {
            std::string_view cookie_name = trim_right(val.substr(0, eq));
            
            size_t cookie_val_len = (sc == std::string_view::npos) ?
                                    (val.size() - eq - 1) : (sc - eq - 1);
            std::string_view cookie_val = trim_left(
                val.substr(eq + 1, cookie_val_len));
            
            if (!cookie_name.empty()) {
              req.cookies[cookie_name] = cookie_val;
            }
          }
          
          if (sc == std::string_view::npos) {
            break;
          }
          val.remove_prefix(sc + 1);
          val = trim_left(val);
        }
      }
    }
    
    if (data.size() < eol + 2) {
      break;
    }
    data.remove_prefix(eol + 2);
  }

  // Parse POST body if Content-Type is application/x-www-form-urlencoded or multipart/form-data
  auto content_type_it = req.headers.find("Content-Type");
  if (content_type_it != req.headers.end()) {
    auto content_type = content_type_it->second;
    
    // Check if this is form data
    if (content_type.find("application/x-www-form-urlencoded") != std::string_view::npos) {
      parse_form_data(req.body, req.query);
    }
    // Check if this is multipart form data (file upload)
    else if (content_type.find("multipart/form-data") != std::string_view::npos) {
      // Extract boundary from Content-Type header
      // Format: multipart/form-data; boundary=----WebKitFormBoundary...
      auto boundary_pos = content_type.find("boundary=");
      if (boundary_pos != std::string_view::npos) {
        auto boundary_start = boundary_pos + 9;
        auto boundary_end = content_type.find(";", boundary_start);
        if (boundary_end == std::string_view::npos) {
          boundary_end = content_type.size();
        }
        
        // Remove quotes if present
        if (content_type[boundary_start] == '"') {
          boundary_start++;
          boundary_end--;
        }
        
        auto boundary = content_type.substr(boundary_start, boundary_end - boundary_start);
        parse_multipart_form(req.body, boundary, req.query, req.files);
      }
    }
    // Additional content types:
    // - application/json (already stored in body)
    // - text/plain (already stored in body)
  }

  return req;
}

static Agai::HttpMethod parse_method(std::string_view m) {
  switch (m.size()) {
  case 3:
    if (m == "GET")
      return Agai::HttpMethod::GET;
    if (m == "PUT")
      return Agai::HttpMethod::PUT;
    break;

  case 4:
    if (m == "POST")
      return Agai::HttpMethod::POST;
    if (m == "HEAD")
      return Agai::HttpMethod::HEAD;
    break;

  case 5:
    if (m == "PATCH")
      return Agai::HttpMethod::PATCH;
    break;

  case 6:
    if (m == "DELETE")
      return Agai::HttpMethod::DELETE;
    break;

  case 7:
    if (m == "OPTIONS")
      return Agai::HttpMethod::OPTIONS;
    break;
  default:
    return Agai::HttpMethod::UNKNOWN;
  }
  return Agai::HttpMethod::UNKNOWN;
}

#endif