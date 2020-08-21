package buffalox

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gobuffalo/buffalo/render"
)

type inlineRenderer struct {
	ctx    context.Context
	name   string
	reader io.Reader
}

func (r inlineRenderer) ContentType() string {
	ext := filepath.Ext(r.name)
	t := mime.TypeByExtension(ext)
	if t == "" {
		t = "application/octet-stream"
	}

	return t
}

func (r inlineRenderer) Render(w io.Writer, d render.Data) error {
	written, err := io.Copy(w, r.reader)
	if err != nil {
		return err
	}

	ctx, ok := r.ctx.(responsible)
	if !ok {
		return fmt.Errorf("context has no response writer")
	}

	header := ctx.Response().Header()
	disposition := fmt.Sprintf("inline; filename=%s", r.name)
	header.Add("Content-Disposition", disposition)
	contentLength := strconv.Itoa(int(written))
	header.Add("Content-Length", contentLength)

	return nil
}

// Inline renders a file inline automatically setting following headers:
//
//   Content-Type
//   Content-Length
//   Content-Disposition
//
// Content-Type is set using mime#TypeByExtension with the filename's extension. Content-Type will default to
// application/octet-stream if using a filename with an unknown extension.
func Inline(ctx context.Context, name string, r io.Reader) render.Renderer {
	return inlineRenderer{
		ctx:    ctx,
		name:   name,
		reader: r,
	}
}

// Inline renders a file inline automatically setting following headers:
//
//   Content-Type
//   Content-Length
//   Content-Disposition
//
// Content-Type is set using mime#TypeByExtension with the filename's extension. Content-Type will default to
// application/octet-stream if using a filename with an unknown extension.
// func (e *render.Engine) Inline(ctx context.Context, name string, r io.Reader) render.Renderer {
// 	return Inline(ctx, name, r)
// }

type responsible interface {
	Response() http.ResponseWriter
}
