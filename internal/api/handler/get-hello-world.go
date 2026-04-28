package handler

import (
	"net/http"

	"github.com/tommaso-sgroi/fantastic-gin-dealcoholized/internal/api/reqcontext"
)

// GetHelloWorld is an example handler that returns "Hello World!" as plain text.
func (h *Handler) GetHelloWorld(w http.ResponseWriter, _ *http.Request, _ PathParameters, ctx *reqcontext.RequestContext) {
	_, _ = w.Write([]byte("Hello World!"))
}
