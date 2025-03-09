package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"
)

// Global Variables in File Scope
var (
	templateRecords map[string]templates // keep the reocrd of all the templated which are already templated
)

type (
	templates struct {
		Uri          string            // path of the template file
		LastModified time.Time         // date when the file last modified
		Data         template.Template // template data of the file before modified
	}
)

type RenderEngine struct {
	view      []string
	viewCount int
	W         http.ResponseWriter
}

type RenderData map[string]interface{}

func NewRenderHandlerObj(_w http.ResponseWriter) RenderEngine {
	return RenderEngine{
		view:      make([]string, 0),
		viewCount: 0,
		W:         _w,
	}
}

func (rh *RenderEngine) Render(massages string) {
	rh.view = append(rh.view, massages)
	rh.viewCount++
}

func (rh *RenderEngine) StartRender() {
	for i := 0; i < rh.viewCount; i++ {
		// fmt.Println("Rendering :", rh.view[i])
		fmt.Fprint(rh.W, rh.view[i])
	}
	// fmt.Fprint(W, view)
	rh.view = []string{}
	rh.viewCount = 0 //  Reseting view Index
}

func (rh *RenderEngine) RenderView(view func(RenderData) string, renderData RenderData) {
	rh.W.Write([]byte(view(renderData)))
}

/*
 * This function will render go default Html templating tool
 * as argument it will take String to render and data which need to be parsed
 */
func (rh *RenderEngine) RenderTemplate(uri string, data any) error {

	var err error
	// var _html_template *template.Template
	var info os.FileInfo

	// _template, isPresent := templateRecords[uri]
	WriteConsole("Rendering Template\n")
	info, err = os.Stat(uri)
	if err != nil {
		return err
	}
	WriteConsole(info.ModTime().String())

	// if !isPresent { // template is not created already then we will update that in reocrd
	// 	if _html_template, err = template.New("").Parse(ReadFromFile(uri)); err == nil {
	// 		templateRecords[uri] = templates{
	// 			Uri:          uri,
	// 			LastModified: info.ModTime(),
	// 			Data:         *_html_template,
	// 		}
	// 	} else {
	// 		return err
	// 	}
	// } else if _template.LastModified.Compare(info.ModTime()) != 0 { // template already present do other stupid stuff

	// 	if _html_template, err = template.New("").Parse(ReadFromFile(uri)); err == nil {
	// 		_template.LastModified = info.ModTime()
	// 		_template.Data = *_html_template
	// 	} else {
	// 		return err
	// 	}
	// }

	// _template.Data.Execute(rh.W, data)
	// _html_template = nil

	return nil
}
