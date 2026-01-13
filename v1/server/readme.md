# HTTP Request Parsing Logic - Complete Documentation

## Table of Contents
1. [Overview](#overview)
2. [Core Concepts](#core-concepts)
3. [Utility Functions](#utility-functions)
4. [HTTP Method Parsing](#http-method-parsing)
5. [Query String Parsing](#query-string-parsing)
6. [Form Data Parsing](#form-data-parsing)
7. [Multipart Form Data Parsing (File Uploads)](#multipart-form-data-parsing-file-uploads)
8. [Complete HTTP Request Parsing](#complete-http-request-parsing)
9. [Examples and Use Cases](#examples-and-use-cases)

---

## Overview

The HTTP request parser (`persing.cpp`) is a crucial component that converts raw network data (bytes) into structured HTTP request objects. When a client sends an HTTP request to your server, it arrives as a continuous stream of text data. This module's job is to understand that raw data and organize it into a meaningful structure that your application can work with.

### Why HTTP Parsing Matters

When a web browser or API client communicates with a web server, they use the HTTP (HyperText Transfer Protocol) standard. This standard defines a specific format for how requests should be sent:

```
GET /search?q=hello HTTP/1.1
Host: example.com
User-Agent: Chrome/91.0
Cookie: sessionId=abc123

```

The parser reads this exact format and extracts:
- **Method**: What action to perform (GET, POST, DELETE, etc.)
- **Path**: Which resource is being requested (/search)
- **Query Parameters**: Additional data in the URL (?q=hello)
- **Headers**: Metadata about the request (Host, User-Agent, etc.)
- **Cookies**: Session tracking information
- **Body**: The actual data being sent (for POST/PUT requests)

---

## Core Concepts

### String Views: Efficient Memory Management

Throughout this code, you'll see `std::string_view`. This is a special C++ type that represents a "view" into existing string data without copying it. Think of it like a window into data that already exists in memory.

**Why use string_view?**
- **Performance**: No memory copying needed
- **Safety**: Read-only access prevents accidental modifications
- **Efficiency**: Perfect for parsing, where you're breaking one large string into smaller pieces

**Example:**
```cpp
std::string_view data = "GET /path HTTP/1.1";
std::string_view method = data.substr(0, 3);  // method = "GET"
data.remove_prefix(3);  // Move the starting point forward
```

### The HttpRequest Structure

The result of parsing is an `Agai::HttpRequest` object containing:

```cpp
struct HttpRequest {
    HttpMethod method;                                  // GET, POST, etc.
    std::string_view path;                              // /search
    std::string_view http_version;                      // HTTP/1.1
    std::unordered_map<std::string_view, 
        std::string_view> headers;                      // Key-value headers
    std::unordered_map<std::string_view, 
        std::string_view> cookies;                      // Key-value cookies
    std::unordered_map<std::string_view, 
        std::string_view> query;                        // Query parameters
    std::string_view body;                              // Request body
};
```

---

## Utility Functions

### 1. trim_left() - Remove Leading Whitespace

**Purpose:** Removes spaces and tabs from the beginning of a string.

**How it works:**
```cpp
static inline std::string_view trim_left(std::string_view s) {
  while (!s.empty() && std::isspace(s[0])) {
    s.remove_prefix(1);  // Remove first character
  }
  return s;
}
```

**Logic Flow:**
1. Check if string has data (`!s.empty()`)
2. Check if first character is whitespace (`std::isspace(s[0])`)
3. If both true, remove the first character by moving the start pointer forward
4. Repeat until we hit a non-whitespace character
5. Return the trimmed string

**Example:**
```
Input:  "   value"
Step 1: "  value"  (remove first space)
Step 2: " value"   (remove second space)
Step 3: "value"    (first char is not whitespace, stop)
Output: "value"
```

**Use Case:** HTTP headers often have spaces after the colon, like `Content-Type: application/json`. We trim to get just `application/json`.

---

### 2. trim_right() - Remove Trailing Whitespace

**Purpose:** Removes spaces and tabs from the end of a string.

**How it works:**
```cpp
static inline std::string_view trim_right(std::string_view s) {
  while (!s.empty() && std::isspace(s[s.length() - 1])) {
    s.remove_suffix(1);  // Remove last character
  }
  return s;
}
```

**Logic Flow:**
1. Check if string has data
2. Check if last character is whitespace
3. If both true, remove the last character
4. Repeat until we hit a non-whitespace character
5. Return the trimmed string

**Example:**
```
Input:  "HTTP/1.1   "
Step 1: "HTTP/1.1  "  (remove last space)
Step 2: "HTTP/1.1 "   (remove space)
Step 3: "HTTP/1.1"    (last char is not whitespace, stop)
Output: "HTTP/1.1"
```

---

### 3. iequals() - Case-Insensitive Comparison

**Purpose:** Compares two strings while ignoring uppercase/lowercase differences.

**How it works:**
```cpp
static inline bool iequals(std::string_view a, std::string_view b) {
  if (a.size() != b.size()) return false;  // Different lengths
  for (size_t i = 0; i < a.size(); i++) {
    if (std::tolower(a[i]) != std::tolower(b[i])) 
      return false;  // Characters don't match
  }
  return true;  // All characters match
}
```

**Logic Flow:**
1. First check: Do the strings have the same length? If not, they can't be equal
2. Loop through each character position
3. Convert both characters to lowercase for comparison
4. If any character doesn't match, return `false`
5. If all characters match, return `true`

**Example:**
```
Compare "Cookie" with "COOKIE":
Position 0: tolower('C') vs tolower('C') ✓
Position 1: tolower('o') vs tolower('O') ✓ (both become 'o')
Position 2: tolower('o') vs tolower('O') ✓
Position 3: tolower('k') vs tolower('K') ✓
Position 4: tolower('i') vs tolower('I') ✓
Position 5: tolower('e') vs tolower('E') ✓
Result: true (strings are equal)
```

**Why needed:** HTTP headers are case-insensitive by standard. "Cookie", "cookie", and "COOKIE" should all be treated the same.

---

## HTTP Method Parsing

### parse_method() - Identify the HTTP Action

**Purpose:** Converts text like "GET" or "POST" into a method identifier.

**How it works:**
```cpp
static Agai::HttpMethod parse_method(std::string_view m) {
  switch (m.size()) {
  case 3:
    if (m == "GET") return Agai::HttpMethod::GET;
    if (m == "PUT") return Agai::HttpMethod::PUT;
    break;
  case 4:
    if (m == "POST") return Agai::HttpMethod::POST;
    if (m == "HEAD") return Agai::HttpMethod::HEAD;
    break;
  case 5:
    if (m == "PATCH") return Agai::HttpMethod::PATCH;
    break;
  case 6:
    if (m == "DELETE") return Agai::HttpMethod::DELETE;
    break;
  case 7:
    if (m == "OPTIONS") return Agai::HttpMethod::OPTIONS;
    break;
  default:
    return Agai::HttpMethod::UNKNOWN;
  }
  return Agai::HttpMethod::UNKNOWN;
}
```

**Logic Flow:**

This function uses an optimization strategy called "length-based dispatch":

1. **First Filter by Length**: Check how many characters are in the method string
   - 3 characters: Could be GET or PUT
   - 4 characters: Could be POST or HEAD
   - 5 characters: Must be PATCH
   - 6 characters: Must be DELETE
   - 7 characters: Must be OPTIONS

2. **Then Compare Strings**: Only compare within the correct length group

**Why This Approach?**
- Fast: Checking string length is O(1) - instant
- Efficient: We eliminate 80% of possibilities with one check
- Readable: Each case is clearly organized

**Common HTTP Methods:**

| Method | Purpose | Meaning |
|--------|---------|---------|
| GET | Retrieve | Fetch a resource (no data changes) |
| POST | Create | Submit data to create something |
| PUT | Replace | Update an entire resource |
| PATCH | Modify | Partially update a resource |
| DELETE | Remove | Delete a resource |
| HEAD | Check | Like GET but without the body |
| OPTIONS | Describe | Ask what methods are allowed |

**Example:**
```
Input: "GET"
Step 1: Length is 3
Step 2: Check if == "GET"? YES
Output: HttpMethod::GET
```

---

## Query String Parsing

### parse_query() - Extract URL Parameters

**Purpose:** When a URL has `?key1=value1&key2=value2`, this function breaks it into key-value pairs.

**How it works:**
```cpp
static void parse_query(
            std::string_view q,
            std::unordered_map<std::string_view, 
            std::string_view> &out
          ) {
  while (!q.empty()) {
    auto eq = q.find('=');      // Find the equals sign
    auto amp = q.find('&');     // Find the ampersand
    
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
```

**Logic Flow - Step by Step:**

Let's parse `search=hello&filter=active&sort=date`:

**Iteration 1:**
```
Current data: "search=hello&filter=active&sort=date"
Step 1: Find '=' → position 6
Step 2: Find '&' → position 12
Step 3: Calculate value length: 12 - 6 - 1 = 5
Step 4: Extract key: "search" (positions 0-6)
Step 5: Extract value: "hello" (positions 7-11)
Step 6: Store in map: query["search"] = "hello"
Step 7: Move past the '&': remove first 13 characters
```

**Iteration 2:**
```
Current data: "filter=active&sort=date"
Step 1: Find '=' → position 6
Step 2: Find '&' → position 13
Step 3: Calculate value length: 13 - 6 - 1 = 6
Step 4: Extract key: "filter"
Step 5: Extract value: "active"
Step 6: Store in map: query["filter"] = "active"
Step 7: Move past the '&': remove first 14 characters
```

**Iteration 3:**
```
Current data: "sort=date"
Step 1: Find '=' → position 4
Step 2: Find '&' → NOT FOUND (npos)
Step 3: Value length: 9 - 4 - 1 = 4 (use total length)
Step 4: Extract key: "sort"
Step 5: Extract value: "date"
Step 6: Store in map: query["sort"] = "date"
Step 7: No '&' found, so break
```

**Error Handling:**
```cpp
if (eq == std::string_view::npos) {
  break;  // No '=' sign, malformed pair
}

if (!key.empty()) {
  out[key] = value;  // Skip empty keys
}
```

**Example URL:**
```
Original: /search?q=hello&page=1&limit=10
Path: /search
Query String: q=hello&page=1&limit=10

After parsing:
query["q"] = "hello"
query["page"] = "1"
query["limit"] = "10"
```

---

## Form Data Parsing

### parse_form_data() - Extract POST Form Data

**Purpose:** When a form is submitted via POST with `Content-Type: application/x-www-form-urlencoded`, this function parses the form fields from the request body.

**How it works:**
```cpp
static void parse_form_data(
            std::string_view body,
            std::unordered_map<std::string_view, 
            std::string_view> &out
          ) {
  // Form data uses the same format as query strings
  parse_query(body, out);
}
```

**Logic:**

Form data sent in POST requests uses the exact same format as query strings:
- Key-value pairs separated by `=`
- Multiple pairs separated by `&`
- Format: `fieldname=value&another=value2`

**Example - HTML Form Submission:**

When a user submits an HTML form:
```html
<form method="POST" action="/login">
  <input type="text" name="username">
  <input type="password" name="password">
  <button type="submit">Login</button>
</form>
```

The browser sends:
```
POST /login HTTP/1.1
Host: example.com
Content-Type: application/x-www-form-urlencoded
Content-Length: 26

username=john&password=secret123
```

**After Parsing:**
```
method = POST
path = "/login"
headers["Content-Type"] = "application/x-www-form-urlencoded"
headers["Content-Length"] = "26"
body = "username=john&password=secret123"
form_data["username"] = "john"
form_data["password"] = "secret123"
```

**Common Use Cases:**

1. **Login Form:**
```
Form sends: username=alice&password=pass123&remember=on
Parsed: form_data["username"] = "alice"
        form_data["password"] = "pass123"
        form_data["remember"] = "on"
```

2. **Search Form:**
```
Form sends: q=javascript&category=tutorial&sort=recent
Parsed: form_data["q"] = "javascript"
        form_data["category"] = "tutorial"
        form_data["sort"] = "recent"
```

3. **Contact Form:**
```
Form sends: name=John+Doe&email=john@example.com&message=Hello
Parsed: form_data["name"] = "John+Doe"  (note: URL decoding not yet implemented)
        form_data["email"] = "john@example.com"
        form_data["message"] = "Hello"
```

**How the Server Detects Form Data:**

The parser automatically detects and parses form data by checking the `Content-Type` header:

```cpp
auto content_type_it = req.headers.find("Content-Type");
if (content_type_it != req.headers.end()) {
  auto content_type = content_type_it->second;
  
  // Check if this is form data
  if (content_type.find("application/x-www-form-urlencoded") 
      != std::string_view::npos) {
    parse_form_data(req.body, req.form_data);
  }
}
```

**Key Differences from Query Strings:**

| Aspect | Query String | Form Data |
|--------|--------------|-----------|
| Location | URL | Request body |
| Method | Usually GET | Usually POST, PUT |
| Header | None | `Content-Type: application/x-www-form-urlencoded` |
| Limit | ~2000 characters | Larger (up to Content-Length) |
| Usage | Filtering, pagination | User input, credentials, file data |
| Field Name | `request.query` | `request.form_data` |

---

## Multipart Form Data Parsing (File Uploads)

### parse_multipart_form() - Extract Files and Mixed Form Data

**Purpose:** Parse forms that contain both text fields AND file uploads. This is the standard format for file upload forms in web browsers.

**How it works:**

Multipart form data uses a special boundary-based format to separate different parts of the request:

```
POST /upload HTTP/1.1
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary123

------WebKitFormBoundary123
Content-Disposition: form-data; name="username"

john_doe
------WebKitFormBoundary123
Content-Disposition: form-data; name="avatar"; filename="profile.jpg"
Content-Type: image/jpeg

[binary image data here]
------WebKitFormBoundary123--
```

**Data Structure: UploadedFile**

```cpp
struct UploadedFile {
  std::string_view filename;      // "profile.jpg"
  std::string_view mime_type;     // "image/jpeg"
  std::string_view content;       // Raw binary file data
  std::string_view field_name;    // "avatar" (form field name)
};
```

**How It Works:**

1. **Extract Boundary** from `Content-Type` header:
   ```
   multipart/form-data; boundary=----WebKitFormBoundary123
   Extract: ----WebKitFormBoundary123
   ```

2. **For Each Part**:
   - Find `--{boundary}` marker
   - Read headers (Content-Disposition, Content-Type)
   - Extract field name and filename (if present)
   - Read content until next boundary

3. **Categorize**:
   - **With filename** → File upload → `request.files` vector
   - **Without filename** → Text field → `request.form_fields` map

### Example: Photo Upload Form

**HTML:**
```html
<form method="POST" enctype="multipart/form-data" action="/profile">
  <input type="text" name="bio">
  <input type="file" name="avatar">
  <button type="submit">Upload</button>
</form>
```

**After Parsing:**
```
form_fields["bio"] = "I love photography"

files[0]:
  field_name = "avatar"
  filename = "me.jpg"
  mime_type = "image/jpeg"
  content = [binary JPEG data]
```

### Common File Upload Scenarios

**Single File:**
```cpp
for (const auto &file : req.files) {
  std::cout << "File: " << file.filename << "\n";
  std::cout << "Type: " << file.mime_type << "\n";
  std::cout << "Size: " << file.content.size() << " bytes\n";
}
```

**Multiple Files:**
```cpp
if (req.files.size() > 0) {
  // Process each uploaded file
  for (const auto &file : req.files) {
    // Save to disk or process
    // file.content has the raw binary data
  }
}
```

**Mixed Form (Text + Files):**
```cpp
// Access text fields
auto title = req.form_fields.find("title");

// Access uploaded files
for (const auto &file : req.files) {
  // Process file
}
```

### Supported File Types

Format-agnostic - handles all file types:
- **Images**: PNG, JPG, GIF, WebP, SVG, BMP
- **Documents**: PDF, DOCX, XLSX, TXT, CSV
- **Archives**: ZIP, TAR, RAR, 7Z
- **Media**: MP3, MP4, WebM, AVI
- **Any binary**: Executables, databases, etc.

MIME type is extracted from each part's `Content-Type` header.

---

## Complete HTTP Request Parsing

### parse_request() - The Main Parser

**Purpose:** Converts raw HTTP bytes into a structured HttpRequest object.

**How it works:**

This is the most complex function, so let's break it into sections:

#### Section 1: Parse the Request Line

```cpp
auto line_end = data.find("\r\n");
if (line_end == std::string_view::npos) {
  return req;  // Malformed: no CRLF found
}

auto line = data.substr(0, line_end);
data.remove_prefix(line_end + 2);
```

**What's Happening:**
1. HTTP uses `\r\n` (CRLF) to mark line endings (Windows convention)
2. Find the first line (ends with `\r\n`)
3. Extract it
4. Remove it from data (move start pointer forward)

**Example:**
```
Raw data:
"GET /search HTTP/1.1\r\nHost: example.com\r\n..."

Step 1: Find "\r\n" at position 20
Step 2: Extract "GET /search HTTP/1.1"
Step 3: Remove from data, now pointing at "Host: example.com\r\n..."
```

#### Section 2: Parse Method, Path, and Version

```cpp
auto sp1 = line.find(' ');     // Space 1: between method and path
auto sp2 = line.find(' ', sp1 + 1);  // Space 2: between path and version

req.method = parse_method(line.substr(0, sp1));
req.http_version = trim_right(line.substr(sp2 + 1));

auto path = line.substr(sp1 + 1, sp2 - sp1 - 1);
```

**Example:**
```
Line: "GET /search?q=hello HTTP/1.1"

Finding spaces:
Space 1 position: 3  (after "GET")
Space 2 position: 18 (after "/search?q=hello")

Extraction:
method = "GET" → HttpMethod::GET
path = "/search?q=hello"
http_version = "HTTP/1.1"
```

#### Section 3: Parse Path and Query String

```cpp
auto qpos = path.find('?');

if (qpos != std::string_view::npos) {
  req.path = path.substr(0, qpos);
  parse_query(path.substr(qpos + 1), req.query);
} else {
  req.path = path;
}
```

**Logic:**
1. Look for `?` in the path (query string separator)
2. If found: Split into path and query parts
3. If not found: Entire thing is just the path

**Example:**
```
Input path: "/search?q=hello&page=1"

Step 1: Find '?' at position 7
Step 2: Extract path: "/search"
Step 3: Extract query: "q=hello&page=1"
Step 4: Parse query string into key-value pairs
```

#### Section 4: Parse Headers

```cpp
while (true) {
  auto eol = data.find("\r\n");
  
  if (eol == 0) {
    // Empty line marks end of headers
    data.remove_prefix(2);
    req.body = data;
    break;
  }

  auto h = data.substr(0, eol);
  auto colon = h.find(':');
  
  if (colon != std::string_view::npos && colon + 1 < h.size()) {
    auto key = h.substr(0, colon);
    auto val = trim_left(h.substr(colon + 1));
    
    req.headers[key] = val;
```

**Logic:**
1. Loop through remaining data, line by line
2. Empty line (`\r\n`) signals end of headers, rest is body
3. For each header line, find the colon separator
4. Key is before the colon, value is after
5. Trim whitespace from the value
6. Store in headers map

**Example Raw Data:**
```
"Host: example.com\r\n"
"Content-Type: application/json\r\n"
"Content-Length: 42\r\n"
"\r\n"
"{"message": "hello world"}"

Step 1: Process "Host: example.com"
        headers["Host"] = "example.com"
Step 2: Process "Content-Type: application/json"
        headers["Content-Type"] = "application/json"
Step 3: Process "Content-Length: 42"
        headers["Content-Length"] = "42"
Step 4: Hit empty line, everything after is body
        body = "{"message": "hello world"}"
```

#### Section 5: Parse Cookies

```cpp
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
```

**Logic:**
1. If header name is "Cookie" (case-insensitive check)
2. Parse like query string but with `;` as separator instead of `&`
3. Format: `name1=value1; name2=value2; name3=value3`
4. Extract each cookie into the cookies map

**Example:**
```
Cookie header: "sessionId=abc123; userId=456; theme=dark"

Iteration 1:
  name: "sessionId"
  value: "abc123"
  cookies["sessionId"] = "abc123"

Iteration 2:
  name: "userId"
  value: "456"
  cookies["userId"] = "456"

Iteration 3:
  name: "theme"
  value: "dark"
  cookies["theme"] = "dark"
```

---

## Examples and Use Cases

### Example 1: Simple GET Request

**Raw HTTP Data:**
```
GET /users HTTP/1.1\r\n
Host: api.example.com\r\n
User-Agent: Mozilla/5.0\r\n
\r\n
```

**After Parsing:**
```
method = GET
path = "/users"
http_version = "HTTP/1.1"
headers["Host"] = "api.example.com"
headers["User-Agent"] = "Mozilla/5.0"
body = (empty)
```

---

### Example 2: POST Request with Query and Body

**Raw HTTP Data:**
```
POST /users?admin=true HTTP/1.1\r\n
Host: api.example.com\r\n
Content-Type: application/json\r\n
Content-Length: 26\r\n
\r\n
{"name":"John","age":30}
```

**After Parsing:**
```
method = POST
path = "/users"
http_version = "HTTP/1.1"
query["admin"] = "true"
headers["Host"] = "api.example.com"
headers["Content-Type"] = "application/json"
headers["Content-Length"] = "26"
body = {"name":"John","age":30}
form_data = (empty - JSON body, not form data)
```

---

### Example 3: POST Form Data Submission

**Raw HTTP Data:**
```
POST /login HTTP/1.1\r\n
Host: example.com\r\n
Content-Type: application/x-www-form-urlencoded\r\n
Content-Length: 32\r\n
\r\n
username=alice&password=secret
```

**After Parsing:**
```
method = POST
path = "/login"
http_version = "HTTP/1.1"
headers["Host"] = "example.com"
headers["Content-Type"] = "application/x-www-form-urlencoded"
headers["Content-Length"] = "32"
body = "username=alice&password=secret"
form_data["username"] = "alice"
form_data["password"] = "secret"
```

---

### Example 4: Request with Cookies

**Raw HTTP Data:**
```
GET /dashboard HTTP/1.1\r\n
Host: app.example.com\r\n
Cookie: sessionId=xyz789; userId=12345; preferences=dark_mode\r\n
\r\n
```

**After Parsing:**
```
method = GET
path = "/dashboard"
headers["Host"] = "app.example.com"
headers["Cookie"] = "sessionId=xyz789; userId=12345; preferences=dark_mode"
cookies["sessionId"] = "xyz789"
cookies["userId"] = "12345"
cookies["preferences"] = "dark_mode"
body = (empty)
```

---

### Example 5: DELETE Request with Complex Path

**Raw HTTP Data:**
```
DELETE /api/posts/42?permanent=true&notify=admin HTTP/1.1\r\n
Host: api.example.com\r\n
Authorization: Bearer token123\r\n
\r\n
```

**After Parsing:**
```
method = DELETE
path = "/api/posts/42"
query["permanent"] = "true"
query["notify"] = "admin"
headers["Host"] = "api.example.com"
headers["Authorization"] = "Bearer token123"
body = (empty)
```

---

### Example 6: File Upload with Multipart Form Data

**Raw HTTP Data:**
```
POST /upload HTTP/1.1\r\n
Host: app.example.com\r\n
Content-Type: multipart/form-data; boundary=----WebKit123\r\n
Content-Length: 1024\r\n
\r\n
------WebKit123\r\n
Content-Disposition: form-data; name="title"\r\n
\r\n
My Photo\r\n
------WebKit123\r\n
Content-Disposition: form-data; name="image"; filename="photo.jpg"\r\n
Content-Type: image/jpeg\r\n
\r\n
[binary JPEG data - 856 bytes]\r\n
------WebKit123--\r\n
```

**After Parsing:**
```
method = POST
path = "/upload"
http_version = "HTTP/1.1"
headers["Host"] = "app.example.com"
headers["Content-Type"] = "multipart/form-data; boundary=----WebKit123"
headers["Content-Length"] = "1024"
body = [full raw multipart data]

form_fields["title"] = "My Photo"

files[0]:
  field_name = "image"
  filename = "photo.jpg"
  mime_type = "image/jpeg"
  content = [856 bytes of JPEG binary data]
```

**Usage in Handler:**
```cpp
Agai::Response handleUpload(const Agai::HttpRequest &req) {
  // Access text field
  auto title_it = req.form_fields.find("title");
  if (title_it != req.form_fields.end()) {
    std::cout << "Title: " << title_it->second << "\n";
  }
  
  // Access uploaded file
  if (!req.files.empty()) {
    const auto &file = req.files[0];
    std::cout << "Filename: " << file.filename << "\n";
    std::cout << "Size: " << file.content.size() << " bytes\n";
    
    // Save file or process it
    // file.content contains raw binary data
  }
  
  return Agai::Response("Upload successful");
}
```

---

### Example 7: Multiple Files with Form Data

**Scenario:** Photo gallery upload with title and description

**After Parsing:**
```
form_fields["title"] = "My Gallery"
form_fields["description"] = "Summer photos"

files[0]:
  field_name = "images"
  filename = "photo1.jpg"
  mime_type = "image/jpeg"
  content = [4KB of image data]

files[1]:
  field_name = "images"
  filename = "photo2.png"
  mime_type = "image/png"
  content = [5KB of image data]

files[2]:
  field_name = "images"
  filename = "photo3.jpg"
  mime_type = "image/jpeg"
  content = [3KB of image data]
```

**Handler Code:**
```cpp
Agai::Response handleGalleryUpload(const Agai::HttpRequest &req) {
  // Get metadata
  std::string title = std::string(req.form_fields.find("title")->second);
  std::string desc = std::string(req.form_fields.find("description")->second);
  
  // Process all files
  for (const auto &file : req.files) {
    std::cout << "Processing: " << file.filename << "\n";
    // Save to disk, resize, thumbnail, etc.
  }
  
  return Agai::Response("Gallery uploaded");
}
```

---

## Key Takeaways

1. **Efficiency**: Uses `string_view` to avoid copying large amounts of data
2. **Robustness**: Handles malformed input gracefully
3. **Standards Compliance**: Follows HTTP/1.1 RFC specifications
4. **Organization**: Separates concerns into focused utility and parsing functions
5. **Performance**: Uses length-based dispatch for fast method identification
6. **File Handling**: Safe, easy file uploads with automatic directory management

This parser and file system is the foundation that allows your web server to understand client requests and handle file uploads securely!
