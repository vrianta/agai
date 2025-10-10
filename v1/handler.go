package agai

import (
	"bytes"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/vrianta/agai/v1/config"
	"github.com/vrianta/agai/v1/internal/template"
	"github.com/vrianta/agai/v1/log"
	"github.com/vrianta/agai/v1/utils"
)

// // Handler processes incoming HTTP requests and manages user sessions.
// // It checks if the user has an existing session and handles session creation or validation.
// // Based on the session and route, it invokes the appropriate controller method.
// // Parameters:
// // - w: The HTTP response writer.
// // - r: The HTTP request.
type (
	controllerInterface interface { // Resembeles Controller Package
		GET() View
		POST() View
		PUT() View
		DELETE() View
		PATCH() View
		HEAD() View
		OPTIONS() View
		Init(w http.ResponseWriter, r *http.Request)
	}
	routes map[string]func() controllerInterface
)

var routeTable routes = make(routes)
var template_bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

func Handler(w http.ResponseWriter, r *http.Request) {

	if _c, found := routeTable[r.URL.Path]; found {

		tempController := _c()
		tempController.Init(w, r)
		// log.WriteLogf("session ID not found\n")
		runRequest(w, r, tempController)

		// if sessionID, err := session.GetSessionID(r); err == nil && sessionID != "" { // it means the user had the session ID
		// 	if sess, _ := session.Get(&sessionID, w, r); sess != nil {
		// 		if tempController, ok := sess.Controller[r.URL.Path]; ok {
		// 			tempController.Init(w, r, sess)
		// 			// log.WriteLogf("controller found in sesion\n")
		// 			runRequest(w, r, tempController)
		// 			return
		// 		} else {
		// 			tempController := _c
		// 			tempController.Init(w, r, sess)
		// 			sess.Controller[r.URL.Path] = tempController
		// 			// log.WriteLogf("controller not found in sesion\n")
		// 			runRequest(w, r, tempController)
		// 			return
		// 		}
		// 	}
		// }

	} else {
		http.Error(w, "404 Error : Route not found ", http.StatusNotFound)
	}
}

func runRequest(w http.ResponseWriter, r *http.Request, c controllerInterface) {

	switch r.Method {
	case "GET":
		vfucn := c.GET()
		view := vfucn() // GET method of controller returns a view
		if view.AsJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			return
		}

		get_template := __template.GET()
		if !config.GetWebConfig().Build {
			// log.WriteLogf("Updating the Template")
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.Response.Get()); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	case "POST":
		vfucn := c.POST()
		view := vfucn() // GET method of controller returns a view
		if view.AsJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			return
		}

		get_template := __template.POST()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.Response.Get()); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	case "DELETE":
		vfucn := c.DELETE()
		view := vfucn() // GET method of controller returns a view
		if view.AsJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			return
		}

		get_template := __template.DELETE()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.Response.Get()); err != nil {
			log.Error("Error rendering template: %T", err)
			return
		}
	case "PATCH":
		vfucn := c.PATCH()
		view := vfucn() // GET method of controller returns a view
		if view.AsJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			return
		}

		get_template := __template.PATCH()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.Response.Get()); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	case "PUT":
		vfucn := c.PUT()
		view := vfucn() // GET method of controller returns a view
		if view.AsJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			return
		}

		get_template := __template.PUT()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.Response.Get()); err != nil {
			panic(err)
		}
	case "HEAD":
		vfucn := c.HEAD()
		view := vfucn() // GET method of controller returns a view
		if view.AsJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			return
		}

		get_template := __template.HEAD()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.Response.Get()); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	case "OPTIONS":
		vfucn := c.OPTIONS()
		view := vfucn() // GET method of controller returns a view
		if view.AsJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			return
		}

		get_template := __template.OPTIONS()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.Response.Get()); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	default:
		log.WriteLogf("Getting Default Method")
		vfucn := c.GET()
		view := vfucn() // GET method of controller returns a view
		if view.AsJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			return
		}

		get_template := __template.INDEX()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.Response.Get()); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	}
}

func executeTemplate(w http.ResponseWriter, _template *template.Context, __response any) error {
	// Use buffer pool for rendering
	buf := template_bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer template_bufPool.Put(buf)

	switch _template.ViewType {
	case template.ViewTypes.PhpTemplate:
		if _template.Php != nil {
			if err := _template.Php.Execute(buf, __response); err != nil {
				return err
			}
		} else {
			panic("php Template is not registered")
		}
	case template.ViewTypes.HtmlTemplate:
		if _template.Html != nil {
			if err := _template.Html.Execute(buf, __response); err != nil {
				return err
			}
		}
	default:
		if _template.Html != nil {
			if err := _template.Html.Execute(buf, __response); err != nil {
				return err
			}
		}
	}

	w.Write(buf.Bytes())
	return nil
}

type (
	fileCacheEntry struct {
		Uri          string    // path of the template file
		LastModified time.Time // date when the file last modified
		Data         string    // template data of the file before modified
	}
)

var fileCache sync.Map // map[string]FileInfo

// StaticFileHandler serves static files with caching support.
// It checks if the file exists in the cache and serves it directly if the cache is valid.
// Otherwise, it reads the file from disk, caches it, and serves it.
// Parameters:
// - contentType: The MIME type of the file being served.
// Returns:
// - http.HandlerFunc: A handler function for serving static files.
func staticFileHandler(contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_filePath := "." + r.URL.Path

		// Attempt to load from cache
		val, _ := fileCache.Load(_filePath)
		cached, ok := val.(fileCacheEntry)

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
		newRecord := fileCacheEntry{
			Uri:          _filePath,
			LastModified: info.ModTime(),
			Data:         _fileData,
		}
		fileCache.Store(_filePath, newRecord)
		w.Write([]byte(_fileData))
	}
}
