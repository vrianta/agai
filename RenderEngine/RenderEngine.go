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
		templateRecordsMutex.RLock()
		templateRecord, ok := templateRecords[uri]
		templateRecordsMutex.RUnlock()
		if ok {
			return rh.ExecuteTemplate(&templateRecord, templateData)
		} else {
			return fmt.Errorf("template %s not found in records", uri)
		}
	}

	var _html_template *template.Template
	full_template_path := Config.ViewFolder + "/" + uri

	templateRecordsMutex.RLock()
	_template, isPresent := templateRecords[uri]
	templateRecordsMutex.RUnlock()

	if info, err := os.Stat(full_template_path); err != nil {
		return err
	} else {
		if !isPresent {
			return fmt.Errorf("template %s not found in records", uri)
		}
		if _template.LastModified.Compare(info.ModTime()) != 0 {
			if _html_template, err = template.New(uri).Parse(PHPToGoTemplate(Utils.ReadFromFile(full_template_path))); err == nil {
				_template.LastModified = info.ModTime()
				_template.Data = _html_template
				templateRecordsMutex.Lock()
				templateRecords[uri] = _template
				templateRecordsMutex.Unlock()
			} else {
				return err
			}
		}
	}

	return rh.ExecuteTemplate(&_template, templateData)
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
