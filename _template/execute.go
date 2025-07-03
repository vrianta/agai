package template

import (
	"bytes"
	"io"
)

func (t *Struct) Execute(_w io.Writer, __response *Response) error {
	// Use buffer pool for rendering
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	switch t.viewType {
	case viewTypes.phpTemplate:
		if t.php != nil {
			if err := t.php.Execute(buf, *__response); err != nil {
				return err
			}
		}
	case viewTypes.htmlTemplate:
		if t.html != nil {
			if err := t.html.Execute(buf, *__response); err != nil {
				return err
			}
		}
	default:
		if t.html != nil {
			if err := t.html.Execute(buf, *__response); err != nil {
				return err
			}
		}
	}

	_w.Write(buf.Bytes())
	return nil
}
