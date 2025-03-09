package server

import (
	"fmt"
	"net/http"
)

type RenderEngine struct {
	view      []string
	viewCount int
	W         http.ResponseWriter
}

type RenderData map[string]interface{}

func NewRenderHandlerObj(_w http.ResponseWriter) RenderEngine {
	return RenderEngine{
		view:      make([]string, 0),
		viewCount: 0,
		W:         _w,
	}
}

func (rh *RenderEngine) Render(massages string) {
	rh.view = append(rh.view, massages)
	rh.viewCount++
}

func (rh *RenderEngine) StartRender() {
	for i := 0; i < rh.viewCount; i++ {
		// fmt.Println("Rendering :", rh.view[i])
		fmt.Fprint(rh.W, rh.view[i])
	}
	// fmt.Fprint(W, view)
	rh.view = []string{}
	rh.viewCount = 0 //  Reseting view Index
}

func (rh *RenderEngine) RenderView(view func(RenderData) string, renderData RenderData) {
	rh.W.Write([]byte(view(renderData)))
}
