package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/zain2323/cronium/services/fileservice/files"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

type FileHandler struct {
	store  files.Storage
	logger *log.Logger
}

func NewFileHandler(store files.Storage, logger *log.Logger) *FileHandler {
	return &FileHandler{store, logger}
}

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Handling file upload")
	filename := chi.URLParam(r, "filename")
	resourceId := chi.URLParam(r, "resourceId")

	h.logger.Println("Resource Id: ", resourceId)
	h.logger.Println("Filename: ", filename)

	// saving file to the storage system in use
	h.saveFile(resourceId, filename, w, r.Body)
}

func (h *FileHandler) saveFile(resourceId string, fileName string, rw http.ResponseWriter, r io.ReadCloser) {
	h.logger.Println("Saving file to storage for resource id: ", resourceId)

	fp := filepath.Join(resourceId, fileName)
	err := h.store.Save(fp, r)
	if err != nil {
		h.logger.Println("Unable to save file to the storage system. ", err)
		http.Error(rw, "Unable to save file to the storage system", http.StatusInternalServerError)
	}
}
