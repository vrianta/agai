package agai

type (
	Response map[string]any
)

/**
 * @param
 * name : name of the View where you want to send the respnse
 * Data you want to pass to the view
 **/
func (c *Controller) View(name string, data any) View {
	return func() *view {
		return &view{
			name:     name,
			response: data,
		}
	}

}

/**
 * If you want to send the response as json
**/
func (c *Controller) ViewAsJson(data Response) View {
	if c.session != nil {
		data["token"] = c.session.ID
	}

	return func() *view {
		return &view{
			asJson:   true,
			response: data,
		}
	}
}

func (r *Response) Get() any {
	return r
}

func (c *Controller) EmptyResponse() *Response {
	return &Response{}
}
