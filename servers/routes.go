package servers

import (
	"encoding/json"
	"log"
	"net/http"
	"sscm/pkg/certs"
	"strconv"

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
		api.Get("/root", h.SearchCert(certs.RootOnly))
		api.Post("/root", h.CreateRoot)
		api.Get("/root/:id", h.DownloadRoot)
		api.Get("/cert", h.SearchCert(certs.NonRootOnly)) // ?parent=1
		api.Post("/cert", h.CreateCert)
		api.Get("/cert/:id", h.DownloadCert)
		api.Get("/search", h.SearchCert(certs.AllCertificates)) // ?q=
	})
	return mux
}

type Handler struct {
	cm certs.Manager
}

func (h Handler) CreateRoot(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) DownloadRoot(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	crt, err := h.cm.Load(id, true)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/x-x509-ca-cert")
	w.Header().Set("Content-Disposition", `attachment; filename="root.crt"`)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	err = certs.WriteCert(crt.CertificateBytes, w)
	if err != nil {
		log.Println(err)
	}
}

func (h Handler) CreateCert(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) DownloadCert(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) SearchCert(mode certs.SearchMode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		certificates, err := h.cm.Search(q, mode)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res := make([]certs.Info, len(certificates))
		for i, v := range certificates {
			res[i] = v.Info
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Println(err)
		}
	}
}
