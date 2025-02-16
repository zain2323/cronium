package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/zain2323/cronium/services/fileservice/files"
	"io"
	"log"
	"net/http"
	"os"
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

func (h *FileHandler) Download(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Handling file download")

	filename := chi.URLParam(r, "filename")
	resourceId := chi.URLParam(r, "resourceId")

	h.logger.Println("Resource Id: ", resourceId)
	h.logger.Println("Filename: ", filename)

	// get the file
	file, err := h.getFile(resourceId, filename, w)
	if err != nil {
		http.Error(w, "Unable to fetch file from the storage system", http.StatusInternalServerError)
		return
	}
	filename = file.Name()

	// set the required headers
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	http.ServeFile(w, r, filename)
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

func (h *FileHandler) getFile(resourceId string, fileName string, rw http.ResponseWriter) (*os.File, error) {
	h.logger.Println("Fetching file from storage for resource id: ", resourceId)
	fp := filepath.Join(resourceId, fileName)
	file, err := h.store.Get(fp)
	if err != nil {
		h.logger.Println("Unable to fetch file from the storage system. ", err)
		return nil, err
	}
	return file, nil
}
