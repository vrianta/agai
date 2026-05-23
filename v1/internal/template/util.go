package template

import (
	"path/filepath"
	"strings"

	"github.com/vrianta/agai/v1/utils"
)

type templateInfo struct {
	fileType     string
	fileName     string
	folderPath   string
	fileFullPath string
	Uri          string
}

func GetFileData(full_file_name string) templateInfo {

	file_type := strings.TrimPrefix(filepath.Ext(full_file_name), ".") // File extension/type
	file_name := full_file_name[:len(full_file_name)-len(file_type)-1] // Name without extension
	folder_path := utils.JoinPath(view_folder, full_file_name)

	return templateInfo{
		fileType:   file_type,
		fileName:   file_name,
		folderPath: folder_path,
	}
}
