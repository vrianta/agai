package requestHandler

import (
	"bytes"
	"net/http"
	"sync"

	"github.com/vrianta/agai/v1/config"
	"github.com/vrianta/agai/v1/internal/session"
	"github.com/vrianta/agai/v1/internal/template"
	"github.com/vrianta/agai/v1/log"
	"github.com/vrianta/agai/v1/view"
)

// // Handler processes incoming HTTP requests and manages user sessions.
// // It checks if the user has an existing session and handles session creation or validation.
// // Based on the session and route, it invokes the appropriate controller method.
// // Parameters:
// // - w: The HTTP response writer.
// // - r: The HTTP request.
type (
	ControllerInterface interface { // Resembeles Controller Package
		GET() func() view.Context
		POST() func() view.Context
		PUT() func() view.Context
		DELETE() func() view.Context
		PATCH() func() view.Context
		HEAD() func() view.Context
		OPTIONS() func() view.Context
		Init(w http.ResponseWriter, r *http.Request, seesion *session.Instance)
	}

	routes map[string]ControllerInterface
)

var RouteTable routes = make(routes)
var template_bufPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func Handler(w http.ResponseWriter, r *http.Request) {

	if _c, found := RouteTable[r.URL.Path]; found {
		if sessionID, err := session.GetSessionID(r); err == nil && sessionID != "" { // it means the user had the session ID
			if sess, _ := session.Get(&sessionID, w, r); sess != nil {
				if tempController, ok := sess.Controller[r.URL.Path]; ok {
					tempController.Init(w, r, sess)
					// log.WriteLogf("controller found in sesion\n")
					runRequest(w, r, tempController)
					return
				} else {
					tempController := _c
					tempController.Init(w, r, sess)
					sess.Controller[r.URL.Path] = tempController
					// log.WriteLogf("controller not found in sesion\n")
					runRequest(w, r, tempController)
					return
				}
			}
		}
		tempController := _c
		tempController.Init(w, r, nil)
		// log.WriteLogf("session ID not found\n")
		runRequest(w, r, tempController)

	} else {
		http.Error(w, "404 Error : Route not found ", http.StatusNotFound)
	}
}

func runRequest(w http.ResponseWriter, r *http.Request, c ControllerInterface) {

	switch r.Method {
	case "GET":
		vfucn := c.GET()
		view := vfucn() // GET method of controller returns a view
		if view.AsJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			log.Debug("Template is nil for controller %s, no template to execute\n", view.Name)
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			panic("")
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
			log.Debug("Template is nil for controller %s, no template to execute\n", view.Name)
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			panic("")
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
			log.Debug("Template is nil for controller %s, no template to execute\n", view.Name)
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			panic("")
		}

		get_template := __template.DELETE()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.Response.Get()); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	case "PATCH":
		vfucn := c.PATCH()
		view := vfucn() // GET method of controller returns a view
		if view.AsJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			log.Debug("Template is nil for controller %s, no template to execute\n", view.Name)
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			panic("")
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
			log.Debug("Template is nil for controller %s, no template to execute\n", view.Name)
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			panic("")
		}

		get_template := __template.PUT()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := executeTemplate(w, get_template, view.Response.Get()); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	case "HEAD":
		vfucn := c.HEAD()
		view := vfucn() // GET method of controller returns a view
		if view.AsJson {
			// user want the response to be send as json
			w.Write(view.ToJson())
			log.Debug("Template is nil for controller %s, no template to execute\n", view.Name)
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			panic("")
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
			log.Debug("Template is nil for controller %s, no template to execute\n", view.Name)
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			panic("")
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
			log.Debug("Template is nil for controller %s, no template to execute\n", view.Name)
			return
		}
		__template, ok := template.GetTemplate(view.Name)

		if !ok {
			log.Error("No Template Available as such name %s", view.Name)
			panic("")
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
