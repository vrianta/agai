package server

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
)

func NewRenderHandlerObj(_w http.ResponseWriter) RenderEngine {
	return RenderEngine{
		view: make([]byte, 0),
		W:    _w,
	}
}

func (rh *RenderEngine) Render(massages string) {
	rh.view = append(rh.view, []byte(massages)...)
}

func (rh *RenderEngine) StartRender() {
	rh.W.Write(rh.view)
}

func (rh *RenderEngine) RenderGothtml(view func(RenderData) string, renderData RenderData) {
	rh.W.Write([]byte(view(renderData)))
}

func (r *RenderEngine) RenderError(_massage string, _response_code ResponseCode) {
	http.Error(r.W, _massage, int(_response_code))
}

/*
 * This function will render go default Html templating tool
 * as argument it will take String to render and data which need to be parsed
 */
func (rh *RenderEngine) RenderTemplate(uri string, templateData any) error {

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
		if _html_template, err = template.New("").Parse(ReadFromFile(full_template_path)); err == nil {
			templateRecords[uri] = templates{
				Uri:          full_template_path,
				LastModified: info.ModTime(),
				Data:         *_html_template,
			}
			_template = templateRecords[uri]
		} else {
			return err
		}
	} else if _template.LastModified.Compare(info.ModTime()) != 0 { // template already present do other stupid stuff
		if _html_template, err = template.New("").Parse(ReadFromFile(full_template_path)); err == nil {
			_template.LastModified = info.ModTime()
			_template.Data = *_html_template
		} else {
			return err
		}
	}

	var buf bytes.Buffer
	if err = _template.Data.Execute(&buf, templateData); err != nil {
		return err
	}
	rh.view = append(rh.view, buf.Bytes()...)
	_html_template = nil

	return nil
}
