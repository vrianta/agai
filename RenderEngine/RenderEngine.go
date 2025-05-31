package RenderEngine

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/vrianta/Server/Config"
	"github.com/vrianta/Server/Response"
	"github.com/vrianta/Server/Template"
	"github.com/vrianta/Server/Utils"
)

func New(_w http.ResponseWriter) Struct {
	return Struct{
		view: make([]byte, 0),
		W:    _w,
	}
}

func (rh *Struct) Render(massages string) {
	rh.view = append(rh.view, []byte(massages)...)
}

func (rh *Struct) StartRender() {
	rh.W.Write(rh.view)
}

func (r *Struct) RenderError(_massage string, _response_code Response.Code) {
	http.Error(r.W, _massage, int(_response_code))
}

/*
 * This function will render go default Html templating tool
 * as argument it will take String to render and data which need to be parsed
 */
func (rh *Struct) RenderTemplate(uri string, templateData *Template.Response) error {

	if Config.Build {
		if templateRecord, ok := templateRecords[uri]; ok {
			err := rh.ExecuteTemplate(&templateRecord, templateData)
			return err
		} else {
			return fmt.Errorf("template %s not found in records", uri)
		}
	}

	var _html_template *template.Template

	full_template_path := Config.ViewFolder + "/" + uri

	_template, isPresent := templateRecords[uri]
	if info, err := os.Stat(full_template_path); err != nil {
		return err
	} else {
		if !isPresent { // template is not created already then we will update that in reocrd
			return fmt.Errorf("template %s not found in records", uri)
		}
		if _template.LastModified.Compare(info.ModTime()) != 0 { // template already present do other stupid stuff
			if _html_template, err = template.New(uri).Parse(PHPToGoTemplate(Utils.ReadFromFile(full_template_path))); err == nil {
				_template.LastModified = info.ModTime()
				_template.Data = *_html_template
			} else {
				return err
			}
		}
	}

	templateExecuteErr := rh.ExecuteTemplate(&_template, templateData)
	return templateExecuteErr

}

func (rh *Struct) ExecuteTemplate(_template *Template.Struct, templateData *Template.Response) error {

	// Use buffer pool for rendering
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	if err := _template.Data.Execute(rh.W, *templateData); err != nil {
		return err
	}

	rh.W.Write(buf.Bytes())
	return nil
}
