package c2ctrl

import (
	"api/config"
	"api/models"
	"api/repos/filesRepo"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DownloadFile serves files to agents who were ordered to download them.
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	fileID := r.PathValue("fileID")

	file, err := filesRepo.Get(fileID)
	if err != nil {
		log.Println("c2: error: DownloadFileOrder.filesRepo.Get:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if file.Order != models.ORDER_DOWNLOAD {
		log.Println("c2: error: DownloadFileOrder: file order mismatch")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if file.CompletedAt != 0 {
		log.Println("c2: error: DownloadFileOrder: file seems already downloaded")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(*config.DownloadsDirectoryPath, file.OriginalPath)

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("c2: error: DownloadFileOrder.os.ReadFile:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(fileContent); err != nil {
		log.Println("c2: error: DownloadFileOrder.w.Write:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if err := filesRepo.SetCompleted(file.ID); err != nil {
		log.Println("c2: error: DownloadFileOrder.filesRepo.SetCompleted:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// UploadFile handles file uploads by agents who were ordered to upload them.
func UploadFile(w http.ResponseWriter, r *http.Request) {
	fileID := r.PathValue("fileID")

	file, err := filesRepo.Get(fileID)
	if err != nil {
		log.Println("c2: error: UploadFile.filesRepo.Get:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if file.Order != models.ORDER_UPLOAD {
		log.Println("c2: error: DownloadFileOrder: file order mismatch")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if file.CompletedAt != 0 {
		log.Println("c2: error: UploadFile: file seems already uploaded")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	fileContent, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("c2: error: UploadFile.io.ReadAll:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	originalFileName := getFileName(file.OriginalPath)
	if len(originalFileName) == 0 {
		log.Println("c2: error: UploadFile.filepath.Base: file name seems empty")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	newOnDiskName := file.ID + "_" + originalFileName

	if err := os.WriteFile(filepath.Join(*config.UploadsDirectoryPath, newOnDiskName), fileContent, 0777); err != nil {
		log.Println("c2: error: UploadFile.os.WriteFile:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if err := filesRepo.SetCompleted(file.ID); err != nil {
		log.Println("c2: error: UploadFile.filesRepo.SetCompleted:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// getFileName returns file's name from a path.
func getFileName(path string) string {
	// Replace all backslashes with forward slashes.
	normalizedPath := strings.ReplaceAll(path, "\\", "/")
	// Clean the path to handle relative components.
	cleanPath := filepath.Clean(normalizedPath)
	// Extract the base name.
	fileName := filepath.Base(cleanPath)
	// If the path still contains a drive letter (e.g., "C:"), take the last segment after the last "/".
	if idx := strings.Index(fileName, ":"); idx != -1 {
		fileName = filepath.Base(strings.TrimPrefix(cleanPath, fileName[:idx+1]))
	}
	return fileName
}
