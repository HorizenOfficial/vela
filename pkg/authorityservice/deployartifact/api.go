package deployartifact

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/HorizenOfficial/vela/pkg/logger"
)

const (
	defaultMultipartMemBytes = 32 << 20 // 32MB
	bytesInMB                = 1024 * 1024
)

type API struct {
	store          *Store
	maxUploadBytes int64
	log            logger.Logger
}

func NewAPI(store *Store, maxUploadMB int64, log logger.Logger) *API {
	var maxUploadBytes int64
	if maxUploadMB > 0 {
		maxUploadBytes = maxUploadMB * bytesInMB
	}

	return &API{
		store:          store,
		maxUploadBytes: maxUploadBytes,
		log:            log,
	}
}

func (a *API) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if a.maxUploadBytes > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, a.maxUploadBytes)
	}

	if err := r.ParseMultipartForm(defaultMultipartMemBytes); err != nil {
		if isBodyTooLarge(err) {
			http.Error(w, "payload too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}

	file, _, err := r.FormFile("wasm")
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	resp, err := a.store.SaveWASM(file)
	if err != nil {
		if errors.Is(err, ErrEmptyWASM) {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		a.log.Error("deploy/upload: failed to store artifact: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		a.log.Error("deploy/upload: failed to encode response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(append(respBytes, '\n')); err != nil {
		a.log.Error("deploy/upload: failed to write response: %v", err)
	}
}

func isBodyTooLarge(err error) bool {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		return true
	}

	return strings.Contains(err.Error(), "request body too large")
}
