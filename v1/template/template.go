package template

import (
	"bytes"
	htmltemplate "html/template"
	"os"
	"sync"
	"time"

	"github.com/vrianta/agai/v1/config"
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

	content := Utils.ReadFromFile(full_path)

	if !config.GetBuild() {
		// Feature: Adding a javascript to impliment hot reload
		content += `
<script>
const source = new EventSource("http://localhost:8888/hot-reload");
async function isCurrentPageAccessible() {
    const url = window.location.origin;
    try {
        const response = await fetch(url, { method: "HEAD", mode: "no-cors" });
        return response.ok || response.type === "opaque"; 
    } catch (err) {
        return false;
    }
}
async function reloadIfAccessible() {
    const accessible = await isCurrentPageAccessible();
    if (accessible) {
        console.log("[LiveReload] Current page accessible, reloading...");
        window.location.reload();
    } else {
        console.warn("[LiveReload] Current page not accessible, skipping reload");
    }
}

source.onmessage = function(event) {
    if (event.data === "reload") {
        reloadIfAccessible();
    }
    if (event.data === "close") {
        window.close();
    }
};

source.onerror = function(err) {
    console.warn("[LiveReload] Disconnected from server", err);
};
</script>
`
	}
	switch file_type {
	case "php", "gophp":
		if _html_template, err := htmltemplate.New(file_name).Funcs(ReponseFuncMaps).Parse(PHPToGoTemplate(content)); err == nil {
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
		if _html_template, err := htmltemplate.New(file_name).Funcs(ReponseFuncMaps).Parse(content); err == nil {
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
		if _html_template, err := htmltemplate.New(file_name).Funcs(ReponseFuncMaps).Parse(content); err == nil {
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

	content := Utils.ReadFromFile(t.uri)

	if !config.GetBuild() {
		content += `<script>
    const source = new EventSource("http://localhost:8888/hot-reload");

    source.onmessage = function(event) {
        console.log(event)
        if (event.data === "reload") {
            console.log("[LiveReload] Reloading page...");
            window.location.reload();
        }
    };

    source.onerror = function(err) {
        console.warn("[LiveReload] Disconnected from server", err);
    };
</script>`
	}

	// Update the Template if needed
	if info, err := os.Stat(t.uri); err != nil {
		return err
	} else {
		if t.lastModified.Compare(info.ModTime()) != 0 {
			if _html_template, template_err := htmltemplate.New(t.name).Funcs(ReponseFuncMaps).Parse(PHPToGoTemplate(content)); template_err == nil {
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
