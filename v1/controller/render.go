package controller

import (
	Config "github.com/vrianta/agai/v1/config"
	Log "github.com/vrianta/agai/v1/log"
	Template "github.com/vrianta/agai/v1/template"
)

/*
ExecuteTemplate renders the given template with the provided response data.
If not in build mode, updates the template before rendering.
Logs and panics on rendering errors.

Parameters:
- __template: pointer to the Template.Struct to render.
- __response: pointer to Template.Response containing data for the template.

Returns:
- error: if updating the template fails (in dev mode).
*/
func (c *Context) ExecuteTemplate(__template *Template.Struct, __response *Template.Response) error {
	if __template == nil {
		if c.templates.View != nil {
			__template = c.templates.View
		} else {
			c.w.Write(__response.AsJson())
			Log.Debug("Template is nil for controller %s, no template to execute\n", c.View)
			return nil
		}
	}

	if !Config.GetWebConfig().Build {
		__template.Update()
		if err := __template.Execute(c.w, __response); err != nil {
			Log.Error("Error rendering template: %T", err)
			panic(err)
		}
		return nil
	}

	if err := __template.Execute(c.w, __response); err != nil {
		Log.Error("rendering template: %T", err)
		return err
	}
	return nil
}
