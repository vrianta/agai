#include "../agai.h"
#include "utils.h"
#include "logs/logs.h"
#include "../config/config.h"

#include <iostream>
#include <fstream>
#include <filesystem>


// ------------------ utils ------------------

static std::vector<std::string> Agai::Utils::split(const std::string &s, char delim) {
  std::vector<std::string> out;
  std::stringstream ss(s);
  std::string item;
  while (std::getline(ss, item, delim))
    out.push_back(item);
  return out;
}

static std::string Agai::Utils::trim(std::string s) {
  while (!s.empty() && isspace(s.front()))
    s.erase(0, 1);
  while (!s.empty() && isspace(s.back()))
    s.pop_back();
  return s;
}

// ------------------- File Operations -------------------

/**
 * Save a single uploaded file to a specified directory
 * @param file The UploadedFile to save
 * @param directory The directory to save to (uses StaticFilesDir if empty)
 * @return true if successful, false otherwise
 */
bool Agai::Utils::saveFile(const UploadedFile &file, const std::string &directory) {
  std::string save_dir = directory.empty() ? Agai::GetConfig().StaticFilesDir : directory;
  std::string full_path = save_dir + "/" + std::string(file.filename);
  
  return saveFileToPath(file, full_path);
}

/**
 * Save a file to a specific full path
 * @param file The UploadedFile to save
 * @param fullPath The complete path including filename
 * @return true if successful, false otherwise
 */
bool Agai::Utils::saveFileToPath(const UploadedFile &file, const std::string &fullPath) {
  try {
    // Create directory if it doesn't exist
    auto path = std::filesystem::path(fullPath);
    auto parent_dir = path.parent_path();
    
    if (!parent_dir.empty()) {
      std::filesystem::create_directories(parent_dir);
    }
    
    // Open file in binary mode for writing
    std::ofstream outfile(fullPath, std::ios::binary);
    
    if (!outfile.is_open()) {
      Agai::Utils::logf("Error: Could not open file for writing: %s", fullPath.c_str());
      return false;
    }
    
    // Write file content
    outfile.write(file.content.data(), file.content.size());
    
    if (!outfile.good()) {
      Agai::Utils::logf("Error: Failed to write file content to: %s", fullPath.c_str());
      outfile.close();
      return false;
    }
    
    outfile.close();
    Agai::Utils::logf("File saved successfully: %s (%zu bytes)", fullPath.c_str(), file.content.size());
    return true;
    
  } catch (const std::exception &e) {
    Agai::Utils::logf("Exception while saving file: %s", e.what());
    return false;
  }
}

/**
 * Save multiple uploaded files to a specified directory
 * @param files Vector of UploadedFile to save
 * @param directory The directory to save to (uses StaticFilesDir if empty)
 * @return true if all files saved successfully, false if any failed
 */
bool Agai::Utils::saveFiles(const std::vector<UploadedFile> &files, const std::string &directory) {
  if (files.empty()) {
    return true;
  }
  
  std::string save_dir = directory.empty() ? Agai::GetConfig().StaticFilesDir : directory;
  bool all_successful = true;
  
  for (const auto &file : files) {
    if (!saveFile(file, save_dir)) {
      all_successful = false;
    }
  }
  
  return all_successful;
}
