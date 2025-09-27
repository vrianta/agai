package controller

// var template_bufPool = sync.Pool{
// 	New: func() interface{} { return new(bytes.Buffer) },
// }

// func (c *Context) execute(_template *template.Context, __response *Response) error {
// 	// Use buffer pool for rendering
// 	buf := template_bufPool.Get().(*bytes.Buffer)
// 	buf.Reset()
// 	defer template_bufPool.Put(buf)

// 	switch _template.ViewType {
// 	case template.ViewTypes.PhpTemplate:
// 		if _template.Php != nil {
// 			if err := _template.Php.Execute(buf, *__response); err != nil {
// 				return err
// 			}
// 		} else {
// 			panic("php Template is not registered")
// 		}
// 	case template.ViewTypes.HtmlTemplate:
// 		if _template.Html != nil {
// 			if err := _template.Html.Execute(buf, *__response); err != nil {
// 				return err
// 			}
// 		}
// 	default:
// 		if _template.Html != nil {
// 			if err := _template.Html.Execute(buf, *__response); err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	c.w.Write(buf.Bytes())
// 	return nil
// }
