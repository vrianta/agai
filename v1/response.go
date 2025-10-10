package agai

type (
	Response map[string]any
)

/**
 * @param - name : name of the View where you want to send the respnse
 **/
func (r *Response) AsView(name string) View {
	return func() view {
		return view{
			Name:     name,
			Response: r,
		}
	}

}

func (r *Response) AsJson() View {
	return func() view {
		return view{
			AsJson:   true,
			Response: r,
		}
	}

}

func (r *Response) Get() any {
	return r
}

func EmptyResponse() *Response {
	return &Response{}
}
