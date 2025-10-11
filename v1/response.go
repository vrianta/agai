package agai

type (
	Response map[string]any
)

/**
 * @param - name : name of the View where you want to send the respnse
 **/
func (r *Response) ToView(name string) View {
	return func() view {
		return view{
			name:     name,
			response: r,
		}
	}

}

/**
 * If you want to send the response as json
**/
func (r *Response) AsJson() View {
	return func() view {
		return view{
			asJson:   true,
			response: r,
		}
	}
}

func (r *Response) Get() any {
	return r
}

func (c *Controller) EmptyResponse() *Response {
	return &Response{}
}
