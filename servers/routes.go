package servers

import (
	"net/http"
	"sscm/pkg/certs"

	"github.com/go-chi/chi"
)

/*
### [GET] /root - List of root certificates
### [POST] /root - Create a root certificate
### [GET] /root/id - Download specific root certificate (without the key)
### [GET] /cert?root=xx - List all child certificates
### [POST] /cert - Create a child certificate
### [GET] /cert/id - Download specific child certificate and the key
### [GET] /search?q=xx - Search all certificates by query
*/
func Routes(cm certs.Manager) *chi.Mux {
	h := Handler{cm}
	mux := chi.NewMux()
	mux.Route("/api", func(api chi.Router) {
		api.Get("/root", h.ListRoot)
		api.Post("/root", h.CreateRoot)
		api.Get("/root/:id", h.DownloadRoot)
		api.Get("/cert", h.ListCert)
		api.Post("/cert", h.CreateCert)
		api.Get("/cert/:id", h.DownloadCert)
		api.Get("/search", h.SearchCert)
	})
	return mux
}

type Handler struct {
	cm certs.Manager
}

func (h Handler) ListRoot(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) CreateRoot(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) DownloadRoot(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) ListCert(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) CreateCert(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) DownloadCert(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) SearchCert(w http.ResponseWriter, r *http.Request) {

}
