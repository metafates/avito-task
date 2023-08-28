package swagger

import (
	"embed"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
)

//go:embed assets
var swaggerFS embed.FS

const fsRoot = "assets"

func FS() (fs.FS, error) {
	return fs.Sub(swaggerFS, fsRoot)
}

func HTTPFS() (http.FileSystem, error) {
	fsys, err := FS()
	if err != nil {
		return nil, err
	}

	return http.FS(fsys), nil
}

var _ http.Handler = (*swaggerHandler)(nil)

type swaggerHandler struct {
	fileServer http.Handler
	jsonSpec   []byte
}

// ServeHTTP implements http.Handler.
func (s *swaggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/swagger.json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(s.jsonSpec)))
		w.Write(s.jsonSpec)
	default:
		s.fileServer.ServeHTTP(w, r)
	}
}

func NewHandler(spec *openapi3.T) (http.Handler, error) {
	fSys, err := HTTPFS()
	if err != nil {
		return nil, err
	}
	fileServer := http.FileServer(fSys)

	jsonSpec, err := spec.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return &swaggerHandler{
		fileServer: fileServer,
		jsonSpec:   jsonSpec,
	}, nil
}
