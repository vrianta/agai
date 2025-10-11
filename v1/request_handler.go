package agai

import (
	"bytes"
	"net/http"
	"sync"

	"github.com/vrianta/agai/v1/config"
	"github.com/vrianta/agai/v1/internal/template"
	"github.com/vrianta/agai/v1/log"
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
		init(w http.ResponseWriter, r *http.Request)
		execute(w http.ResponseWriter, r *http.Request) view
	}
	routes map[string]func() controllerInterface
)

var routeTable routes = make(routes)
var template_bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

func Handler(w http.ResponseWriter, r *http.Request) {

	if _c, found := routeTable[r.URL.Path]; found {

		log.WriteLogf("url you Hit, %s\n", r.URL.Path)
		tempController := _c()
		tempController.init(w, r)
		runRequest(w, r, tempController)

	} else {
		http.Error(w, "404 Error : Route not found ", http.StatusNotFound)
	}
}

func runRequest(w http.ResponseWriter, r *http.Request, c controllerInterface) {

	switch r.Method {
	case "GET":
		vfucn := c.GET()
		view := vfucn() // GET method of controller returns a view
		if view == nil {
			return
		}
		if view.asJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			return
		}

		get_template := __template.GET()
		if !config.GetWebConfig().Build {
			// log.WriteLogf("Updating the Template")
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.response.Get()); err != nil {
			log.Error("Error rendering template: %T\n", err)
			panic(err)
		}
	case "POST":
		vfucn := c.POST()
		view := vfucn() // GET method of controller returns a view
		if view == nil {
			return
		}
		if view.asJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			return
		}

		get_template := __template.POST()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.response.Get()); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	case "DELETE":
		vfucn := c.DELETE()
		view := vfucn() // GET method of controller returns a view
		if view == nil {
			return
		}
		if view.asJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			return
		}

		get_template := __template.DELETE()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.response.Get()); err != nil {
			log.Error("Error rendering template: %T", err)
			return
		}
	case "PATCH":
		vfucn := c.PATCH()
		view := vfucn() // GET method of controller returns a view
		if view == nil {
			return
		}
		if view.asJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			return
		}

		get_template := __template.PATCH()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.response.Get()); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	case "PUT":
		vfucn := c.PUT()
		view := vfucn() // GET method of controller returns a view
		if view == nil {
			return
		}
		if view.asJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			return
		}

		get_template := __template.PUT()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.response.Get()); err != nil {
			panic(err)
		}
	case "HEAD":
		vfucn := c.HEAD()
		view := vfucn() // GET method of controller returns a view
		if view == nil {
			return
		}
		if view.asJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			return
		}

		get_template := __template.HEAD()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.response.Get()); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	case "OPTIONS":
		vfucn := c.OPTIONS()
		view := vfucn() // GET method of controller returns a view
		if view == nil {
			return
		}
		if view.asJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			return
		}

		get_template := __template.OPTIONS()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.response.Get()); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	default:
		log.WriteLogf("Getting Default Method")
		vfucn := c.GET()
		view := vfucn() // GET method of controller returns a view
		if view == nil {
			return
		}
		if view.asJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			return
		}

		get_template := __template.INDEX()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.response.Get()); err != nil {
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
