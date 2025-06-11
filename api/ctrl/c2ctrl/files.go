package c2ctrl

import (
	"api/config"
	"api/repos/filesRepo"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	fileID := r.PathValue("fileID")

	file, err := filesRepo.Get(fileID)
	if err != nil {
		log.Println("c2: error: UploadFile.filesRepo.Get:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	files, err := filesRepo.GetMultiple()
	fmt.Println(files)

	if file.UploadedAt != 0 {
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

	if err := os.WriteFile(filepath.Join(*config.DownloadsDirectoryPath, newOnDiskName), fileContent, 0777); err != nil {
		log.Println("c2: error: UploadFile.os.WriteFile:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	filesRepo.SetUploaded(file.ID)
}

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
