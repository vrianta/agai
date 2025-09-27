package controller

import (
	"github.com/vrianta/agai/v1/view"
)

type (
	Response map[string]any
)

/**
 * @param - name : name of the View where you want to send the respnse
 **/
func (r *Response) ToView(name string) View {
	return func() view.Context {
		return view.Context{
			Name:     name,
			Response: r,
		}
	}

}

func (r *Response) AsJson() View {
	return func() view.Context {
		return view.Context{
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
