package servers

import (
	"archive/zip"
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sscm/pkg/certs"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

/*
### [GET] /certs?root=xx - List all child certificates (or root if ?root=0)
### [POST] /cert?root=xx - Create a child certificate (or root if ?parent=0)
### [GET] /cert/id - Get detailed information about certificate
### [GET] /download/id - Download specific certificate (root without the key, child with key)
### [GET] /search?q=xx&root=xx - Search all certificates by query
*/

func Routes(cm certs.Manager) *chi.Mux {
	h := Handler{cm}
	mux := chi.NewMux()
	mux.Route("/api", func(api chi.Router) {
		api.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{"http://localhost:3000"},
		}))
		api.Get("/certs", h.ListCert)   // ?root=int
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
	ID      int    `json:"id"`
	Name    string `json:"name"`
	DNS     string `json:"dns"`
	Expired string `json:"expired"`
}

type ListCertResponse struct {
	Parent CertSimpleInfo   `json:"parent,omitempty"`
	Certs  []CertSimpleInfo `json:"certs"`
}

type CreateRequest struct {
	Name         string `json:"name"`
	Country      string `json:"country"`
	Organization string `json:"organization"`
	IP           string `json:"ip"`
	DNS          string `json:"dns"`
}

// ListCert handle [GET] /certs?root=int - List all child certificates (or root if ?root=0)
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

	var res = ListCertResponse{
		Certs: make([]CertSimpleInfo, 0),
	}

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
			ID:      v.ID,
			Name:    v.CommonName,
			DNS:     strings.Join(v.DNSNames, ","),
			Expired: v.ExpiredAt.String(),
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

	var req CreateRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cert := x509.Certificate{
		Subject: pkix.Name{
			Country:      []string{req.Country},
			Organization: []string{req.Organization},
			CommonName:   req.Name,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(10, 0, 0),
		IPAddresses: []net.IP{},
		DNSNames:    []string{},
	}
	if req.IP != "" {
		cert.IPAddresses = []net.IP{net.ParseIP(req.IP)}
	}
	if req.DNS != "" {
		cert.DNSNames = []string{req.DNS}
	}

	var certByte []byte
	var certKey certs.PrivateKey

	if rootId == 0 {
		certByte, certKey, err = certs.CreateRootCA(&cert)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		// load parent certificate
		rootData, err := h.cm.Load(rootId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(rootData.CertificateBytes) == 0 {
			fmt.Println("empty root bytes")
			return
		}
		rootCert, err := certs.LoadCert(bytes.NewReader(rootData.CertificateBytes))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rootKey, err := certs.LoadKey(bytes.NewReader(rootData.PrivateKeyBytes))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		certByte, certKey, err = certs.CreateServerCertificate(&cert, rootCert, rootKey)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	bufC := bytes.Buffer{}
	if err := certs.WriteCert(certByte, &bufC); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bufK := bytes.Buffer{}
	if err := certs.WriteKey(certKey, &bufK); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := certs.Data{
		ParentID:         rootId,
		CertificateBytes: bufC.Bytes(),
		PrivateKeyBytes:  bufK.Bytes(),
		Info: certs.Info{
			CommonName:   cert.Subject.CommonName,
			Country:      cert.Subject.Country[0],
			Organization: cert.Subject.Organization[0],
			IPAddresses:  []string{req.IP},
			DNSNames:     cert.DNSNames,
		},
	}

	if err := h.cm.Store(&data); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// DownloadCert handle [GET] /download/{id} - Download specific certificate (root without the key, child with key)
func (h Handler) DownloadCert(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
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

	name := slugify(crt.CommonName)

	// Download Root Certificate
	if crt.ParentID == 0 {
		w.Header().Set("Content-Type", "application/x-x509-ca-cert")
		w.Header().Set("Content-Disposition", `attachment; filename="`+name+`.crt"`)
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.WriteHeader(http.StatusOK)
		err = certs.WriteCert(crt.CertificateBytes, w)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// Download Standard Certificate
	buf := bytes.Buffer{}
	zw := zip.NewWriter(&buf)
	cw, err := zw.Create(name + ".crt")
	if err != nil {
		log.Println(err)
	}
	err = certs.WriteCert(crt.CertificateBytes, cw)
	if err != nil {
		log.Println(err)
	}
	kw, err := zw.Create(name + ".key")
	if err != nil {
		log.Println(err)
	}
	_, err = kw.Write(crt.PrivateKeyBytes)
	if err != nil {
		log.Println(err)
	}
	err = zw.Close()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="`+name+`.zip"`)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf.Bytes())
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

func slugify(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "_")
}
