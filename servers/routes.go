package servers

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"log"
	"net/http"
	"sscm/pkg/certs"
	"strconv"

	"github.com/go-chi/chi"
)

/*
### [GET] /cert?root=xx - List all child certificates (or root if ?root=0)
### [POST] /cert?root=xx - Create a child certificate (or root if ?parent=0)
### [GET] /cert/id - Get detailed information about certificate
### [GET] /download/id - Download specific certificate (root without the key, child with key)
### [GET] /search?q=xx&root=xx - Search all certificates by query
*/

func Routes(cm certs.Manager) *chi.Mux {
	h := Handler{cm}
	mux := chi.NewMux()
	mux.Route("/api", func(api chi.Router) {
		api.Get("/cert", h.ListCert)    // ?root=int
		api.Post("/cert", h.CreateCert) // ?root=int
		api.Get("/cert/{id}", h.GetCert)
		api.Get("/download/{id}", h.DownloadCert)
		api.Get("/search", h.SearchCert) // ?q=string&root=int
	})
	return mux
}

type Handler struct {
	cm certs.Manager
}

type CertSimpleInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ListCertResponse struct {
	Parent CertSimpleInfo   `json:"parent,omitempty"`
	Certs  []CertSimpleInfo `json:"certs"`
}

type CreateRequest struct {
}

// ListCert handle [GET] /cert?root=int - List all child certificates (or root if ?root=0)
func (h Handler) ListCert(w http.ResponseWriter, r *http.Request) {
	root := r.URL.Query().Get("root")
	if root == "" {
		root = "0"
	}
	rootId, err := strconv.Atoi(root)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res ListCertResponse

	if rootId > 0 {
		rootData, err := h.cm.Load(rootId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Parent = CertSimpleInfo{
			ID:   rootData.ID,
			Name: rootData.CommonName,
		}
	}

	data, err := h.cm.Search("", rootId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, v := range data {
		res.Certs = append(res.Certs, CertSimpleInfo{
			ID:   v.ID,
			Name: v.CommonName,
		})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// GetCert handle [GET] /cert/{id} - Get detailed information about certificate
func (h Handler) GetCert(w http.ResponseWriter, r *http.Request) {
	sid := chi.URLParam(r, "id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if id < 1 {
		log.Println("Invalid ID:", id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cert, err := h.cm.Load(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_ = cert // TODO convert cert to response
}

// CreateCert handle [POST] /cert?root=int - Create a child certificate (or root if ?parent=0)
func (h Handler) CreateCert(w http.ResponseWriter, r *http.Request) {
	root := r.URL.Query().Get("root")
	if root == "" {
		root = "0"
	}
	rootId, err := strconv.Atoi(root)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cert := x509.Certificate{
		// TODO
	}

	certByte, certKey, err := certs.CreateRootCA(&cert)
	data := certs.Data{
		ParentID:         rootId,
		CertificateBytes: certByte,
		PrivateKeyBytes:  make([]byte, 0),
		Info: certs.Info{
			CommonName:   cert.Subject.CommonName,
			Country:      cert.Subject.Country[0],
			Organization: cert.Subject.Organization[0],
			IPAddresses:  []string{cert.IPAddresses[0].To4().String()},
			DNSNames:     cert.DNSNames,
		},
	}
	certs.WriteKey(certKey, bytes.NewBuffer(data.PrivateKeyBytes))

	if err := h.cm.Store(&data); err != nil {
		// TODO
	}
}

// DownloadCert handle [GET] /download/{id} - Download specific certificate (root without the key, child with key)
func (h Handler) DownloadCert(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	crt, err := h.cm.Load(id)
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

// SearchCert handle [GET] /search?q=str&root=int - Search all certificates by query
func (h Handler) SearchCert(w http.ResponseWriter, r *http.Request) {
	rootId := -1
	if r.URL.Query().Has("root") {
		root := r.URL.Query().Get("root")
		if root == "" {
			root = "0"
		}
		id, err := strconv.Atoi(root)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		rootId = id
	}

	q := r.URL.Query().Get("q")

	certificates, err := h.cm.Search(q, rootId)
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
