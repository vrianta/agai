package router

import (
	"net/http"
	"os"
	"sync"
	"time"

	requesthandler "github.com/vrianta/agai/v1/internal/request_handler"
	"github.com/vrianta/agai/v1/log"
	"github.com/vrianta/agai/v1/utils"
)

type (
	route struct {
		DefaultRoute string
	}

	FileCacheEntry struct {
		Uri          string    // path of the template file
		LastModified time.Time // date when the file last modified
		Data         string    // template data of the file before modified
	}
)

var (
	fileCache sync.Map // map[string]FileInfo
)

/*
 * Create a New Router Object with Default route group example / is the default route for this or /api or /v1 etc
 */
func New(root string) *route {
	return &route{
		DefaultRoute: root,
	}
}

func (r *route) AddRoute(path string, t any) {
	requesthandler.CreateRoute[t](r.DefaultRoute + path)

}

// A Function to Create and Return

// StaticFileHandler serves static files with caching support.
// It checks if the file exists in the cache and serves it directly if the cache is valid.
// Otherwise, it reads the file from disk, caches it, and serves it.
// Parameters:
// - contentType: The MIME type of the file being served.
// Returns:
// - http.HandlerFunc: A handler function for serving static files.
func StaticFileHandler(contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_filePath := "." + r.URL.Path

		// Attempt to load from cache
		val, _ := fileCache.Load(_filePath)
		cached, ok := val.(FileCacheEntry)

		info, err := os.Stat(_filePath)
		if err != nil {
			log.WriteLog(err.Error())
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", contentType)

		// If cached data exists and mod time matches then serve from cache
		if ok && cached.LastModified.Equal(info.ModTime()) {
			w.Write([]byte(cached.Data))
			return
		}

		// Read file from disk and cache it
		_fileData := utils.ReadFromFile(_filePath)
		newRecord := FileCacheEntry{
			Uri:          _filePath,
			LastModified: info.ModTime(),
			Data:         _fileData,
		}
		fileCache.Store(_filePath, newRecord)
		w.Write([]byte(_fileData))
	}
}

// return a list of all the views from the controllers
// loop throgh all the controllers and make a array of strings
// func ListViews() []string {
// 	routerSize := len(routeTable)
// 	if routerSize < 1 {
// 		return nil
// 	}
// 	response := make([]string, 0, routerSize)
// 	for _, controller := range routeTable {
// 		response = append(response, controller.View)
// 	}
// 	return response
// }
