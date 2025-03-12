package server

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
)

// Global Variables in File Scope
var (
	templateRecords = make(map[string]templates) // keep the reocrd of all the templated which are already templated
)

func NewRenderHandlerObj(_w http.ResponseWriter) RenderEngine {
	return RenderEngine{
		view:      make([]byte, 0),
		viewCount: 0,
		W:         _w,
	}
}

func (rh *RenderEngine) Render(massages string) {
	rh.view = append(rh.view, []byte(massages)...)
	rh.viewCount++
}

func (rh *RenderEngine) StartRender() {
	rh.W.Write(rh.view)
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
	var _html_template *template.Template
	var info os.FileInfo

	full_template_path := srvInstance.Config.Views_folder + "/" + uri

	_template, isPresent := templateRecords[uri]
	info, err = os.Stat(full_template_path)
	if err != nil {
		return err
	}

	if !isPresent { // template is not created already then we will update that in reocrd
		// WriteLogf("Trying to fid the View %s/%s", srvInstance.Config.Views_folder, uri)
		if _html_template, err = template.New("").Parse(ReadFromFile(full_template_path)); err == nil {
			templateRecords[uri] = templates{
				Uri:          uri,
				LastModified: info.ModTime(),
				Data:         *_html_template,
			}
			_template = templateRecords[uri]
		} else {
			return err
		}
	} else if _template.LastModified.Compare(info.ModTime()) != 0 { // template already present do other stupid stuff
		// WriteConsole("File Has been Modified")
		if _html_template, err = template.New("").Parse(ReadFromFile(full_template_path)); err == nil {
			_template.LastModified = info.ModTime()
			_template.Data = *_html_template
		} else {
			return err
		}
	}

	// _template.Data.Execute(rh.W, data)
	var buf bytes.Buffer
	err = _template.Data.Execute(&buf, data)
	if err != nil {
		return err
	}
	rh.view = append(rh.view, buf.Bytes()...)
	_html_template = nil

	return nil
}
