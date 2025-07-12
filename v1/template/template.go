package template

import (
	"bytes"
	htmltemplate "html/template"
	"os"
	"sync"
	"time"

	Utils "github.com/vrianta/agai/v1/utils"
)

type (
	Response map[string]any
	ViewType int16

	Struct struct {
		uri          string                 // path of the template file
		name         string                 // name of the file
		lastModified time.Time              // date when the file last modified
		html         *htmltemplate.Template // template data of the file before modified
		php          *htmltemplate.Template // template which will hold the data for php templates
		viewType     ViewType               // type of the view file (html, php, etc.)
	}
)

var (
	EmptyResponse = Response{}
	NoResponse    = Response{}

	viewTypes = struct {
		goTemplate   ViewType
		htmlTemplate ViewType
		phpTemplate  ViewType
	}{
		goTemplate:   0,
		htmlTemplate: 1,
		phpTemplate:  2,
	}

	templateRecordsMutex = &sync.RWMutex{}
	bufPool              = sync.Pool{
		New: func() interface{} { return new(bytes.Buffer) },
	}
)

// Create Template Object stores it in the memory
// name - name of the template
func New(file_path, file_name, file_type string) (*Struct, error) {

	if file_path == "" {
		return nil, nil
	}

	var full_path = file_path + "/" + file_name

	info, err := os.Stat(full_path)
	if err != nil {
		return nil, err
	}

	switch file_type {
	case "php", "gophp":
		if _html_template, err := htmltemplate.New(file_name).Parse(PHPToGoTemplate(Utils.ReadFromFile(full_path))); err == nil {
			return &Struct{
				uri:          full_path,
				name:         file_name,
				lastModified: info.ModTime(),
				php:          _html_template,
				viewType:     viewTypes.phpTemplate,
			}, nil
		} else {
			return nil, err
		}
	case "html", "gohtml":
		if _html_template, err := htmltemplate.New(file_name).Parse(Utils.ReadFromFile(full_path)); err == nil {
			return &Struct{
				uri:          full_path,
				name:         file_name,
				lastModified: info.ModTime(),
				html:         _html_template,
				viewType:     viewTypes.htmlTemplate,
			}, nil
		} else {
			return nil, err
		}
	default:
		if _html_template, err := htmltemplate.New(file_name).Parse(Utils.ReadFromFile(full_path)); err == nil {
			return &Struct{
				uri:          full_path,
				name:         file_name,
				lastModified: info.ModTime(),
				html:         _html_template,
				viewType:     viewTypes.htmlTemplate,
			}, nil
		} else {
			return nil, err
		}
	}
}

/*
uri - full path of the template
*/
func (t *Struct) Update() error {
	templateRecordsMutex.Lock()
	defer templateRecordsMutex.Unlock()

	// Update the Template if needed
	if info, err := os.Stat(t.uri); err != nil {
		return err
	} else {
		if t.lastModified.Compare(info.ModTime()) != 0 {
			if _html_template, template_err := htmltemplate.New(t.name).Parse(PHPToGoTemplate(Utils.ReadFromFile(t.uri))); template_err == nil {
				t.lastModified = info.ModTime()
				switch t.viewType {
				case viewTypes.phpTemplate:
					t.php = _html_template
				case viewTypes.htmlTemplate:
					t.html = _html_template
				default:
					t.html = _html_template
				}
			} else {
				return template_err
			}
		}
	}
	return nil
}
