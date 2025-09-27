package template

import (
	htmltemplate "html/template"
	"os"
	"sync"
	"time"

	"github.com/vrianta/agai/v1/config"
	"github.com/vrianta/agai/v1/utils"
)

type (
	Response map[string]any
	ViewType int16

	Context struct {
		uri          string                 // path of the template file
		name         string                 // name of the file
		lastModified time.Time              // date when the file last modified
		Html         *htmltemplate.Template // template data of the file before modified
		Php          *htmltemplate.Template // template which will hold the data for php templates
		ViewType     ViewType               // type of the view file (html, php, etc.)
		initialised  bool                   // will hold information if the template is initialised
	}

	// holdingh different templates for differect method
	Contexts struct {
		index   *Context // default template store
		get     *Context // Template for GET requests
		post    *Context // Template for POST requests
		delete  *Context // Template for DELETE requests
		patch   *Context // Template for PATCH requests
		put     *Context // Template for PUT requests
		head    *Context // Template for HEAD requests
		options *Context // Template for OPTIONS requests
	}
)

var (
	ViewTypes = struct {
		GoTemplate   ViewType
		HtmlTemplate ViewType
		PhpTemplate  ViewType
	}{
		GoTemplate:   0,
		HtmlTemplate: 1,
		PhpTemplate:  2,
	}

	templateRecordsMutex = &sync.RWMutex{}

	templateRegistry map[string]*Contexts = make(map[string]*Contexts) // holding all the templates in the solution
)

// Create Template Object stores it in the memory
// name - name of the template
func create(file_path, file_name, file_type string) (*Context, error) {

	if file_path == "" {
		return nil, nil
	}

	var full_path = file_path + "/" + file_name

	info, err := os.Stat(full_path)
	if err != nil {
		return nil, err
	}

	content := utils.ReadFromFile(full_path)

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
			return &Context{
				uri:          full_path,
				name:         file_name,
				lastModified: info.ModTime(),
				Php:          _html_template,
				ViewType:     ViewTypes.PhpTemplate,
			}, nil
		} else {
			return nil, err
		}
	case "html", "gohtml":
		if _html_template, err := htmltemplate.New(file_name).Funcs(ReponseFuncMaps).Parse(content); err == nil {
			return &Context{
				uri:          full_path,
				name:         file_name,
				lastModified: info.ModTime(),
				Html:         _html_template,
				ViewType:     ViewTypes.HtmlTemplate,
			}, nil
		} else {
			return nil, err
		}
	default:
		if _html_template, err := htmltemplate.New(file_name).Funcs(ReponseFuncMaps).Parse(content); err == nil {
			return &Context{
				uri:          full_path,
				name:         file_name,
				lastModified: info.ModTime(),
				Html:         _html_template,
				ViewType:     ViewTypes.HtmlTemplate,
			}, nil
		} else {
			return nil, err
		}
	}
}

/*
uri - full path of the template
*/
func (t *Context) Update() error {
	templateRecordsMutex.Lock()
	defer templateRecordsMutex.Unlock()

	content := utils.ReadFromFile(t.uri)

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
				switch t.ViewType {
				case ViewTypes.PhpTemplate:
					t.Php = _html_template
				case ViewTypes.HtmlTemplate:
					t.Html = _html_template
				default:
					t.Html = _html_template
				}
			} else {
				return template_err
			}
		}
	}
	return nil
}

// Returning the Templte Contexts which have all the Method Related Templates
func GetTemplate(name string) (*Contexts, bool) {
	c, ok := templateRegistry[name]

	return c, ok
}

// INDEX returns the index Context handler.
func (c *Contexts) INDEX() *Context {
	return c.index
}

// GET returns the GET Context handler, or falls back to INDEX if nil.
func (c *Contexts) GET() *Context {
	if c.get == nil {
		return c.index
	}
	return c.get
}

// POST returns the POST Context handler, or falls back to INDEX if nil.
func (c *Contexts) POST() *Context {
	if c.post == nil {
		return c.index
	}
	return c.post
}

// PUT returns the PUT Context handler, or falls back to INDEX if nil.
func (c *Contexts) PUT() *Context {
	if c.put == nil {
		return c.index
	}
	return c.put
}

// PATCH returns the PATCH Context handler, or falls back to INDEX if nil.
func (c *Contexts) PATCH() *Context {
	if c.patch == nil {
		return c.index
	}
	return c.patch
}

// DELETE returns the DELETE Context handler, or falls back to INDEX if nil.
func (c *Contexts) DELETE() *Context {
	if c.delete == nil {
		return c.index
	}
	return c.delete
}

// HEAD returns the HEAD Context handler, or falls back to INDEX if nil.
func (c *Contexts) HEAD() *Context {
	if c.head == nil {
		return c.index
	}
	return c.head
}

// OPTIONS returns the OPTIONS Context handler, or falls back to INDEX if nil.
func (c *Contexts) OPTIONS() *Context {
	if c.options == nil {
		return c.index
	}
	return c.options
}

// // TRACE returns the TRACE Context handler, or falls back to INDEX if nil.
// func (c *Contexts) TRACE() *Context {
// 	if c.trace == nil {
// 		return c.index
// 	}
// 	return c.trace
// }

// // CONNECT returns the CONNECT Context handler, or falls back to INDEX if nil.
// func (c *Contexts) CONNECT() *Context {
// 	if c.connect == nil {
// 		return c.index
// 	}
// 	return c.connect
// }
