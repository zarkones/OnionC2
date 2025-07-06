package apictrl

import (
	"api/config"
	"api/core"
	"api/models"
	"api/repos/filesRepo"
	"api/repos/messagesRepo"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gorm.io/gorm"
)

type FileRecord struct {
	Name      string `json:"name"`
	IsDir     bool   `json:"isDir"`
	Timestamp int64  `json:"timestamp"`
}

type RemoteFS struct {
	ID                string       `json:"id"`
	LatestRequestedID string       `json:"latestRequestedId"`
	CurrentDir        string       `json:"currentDir"`
	Content           []FileRecord `json:"content"`
}

func GetRemoteFS(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_FILES_REMOTE_REPO_LIST, nil)
	if reject {
		return
	}

	agentID := r.PathValue("agentID")

	latestFS, err := messagesRepo.GetLatestFS(agentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	latestRequestedFS, err := messagesRepo.GetLatestRequestedFS(agentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var remoteFS RemoteFS

	if err := json.Unmarshal([]byte(latestFS.Response), &remoteFS); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	remoteFS.ID = latestFS.ID
	remoteFS.LatestRequestedID = latestRequestedFS.ID

	if err := json.NewEncoder(w).Encode(&remoteFS); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetUploadedFiles(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_FILES_UPLOADS_REPO_LIST, nil)
	if reject {
		return
	}

	agentID := r.PathValue("agentID")

	fileOrders, err := filesRepo.GetMultipleByAgentID(agentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(fileOrders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	records, err := os.ReadDir(*config.UploadsDirectoryPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(records) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	normalizedFileNames := make([]string, len(fileOrders))
	for i, fileOrder := range fileOrders {
		normalizedFileNames[i] = fileOrder.ID + "_" + core.GetFileName(fileOrder.OriginalPath)
	}

	normalizedRecords := []FileRecord{}

	for _, record := range records {
		if !slices.Contains(normalizedFileNames, record.Name()) {
			continue
		}
		info, err := record.Info()
		if err != nil {
			normalizedRecords = append(normalizedRecords, FileRecord{
				Name:      record.Name(),
				IsDir:     record.IsDir(),
				Timestamp: 0,
			})
			continue
		}
		normalizedRecords = append(normalizedRecords, FileRecord{
			Name:      record.Name(),
			IsDir:     record.IsDir(),
			Timestamp: info.ModTime().UnixNano(),
		})
	}

	if len(normalizedRecords) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(w).Encode(&normalizedRecords); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetDownloadRepositoryFiles(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_FILES_DOWNLOADS_REPO_LIST, nil)
	if reject {
		return
	}

	records, err := os.ReadDir(*config.DownloadsDirectoryPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(records) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	normalizedRecords := make([]FileRecord, len(records))

	for i, record := range records {
		info, err := record.Info()
		if err != nil {
			normalizedRecords[i] = FileRecord{
				Name:      record.Name(),
				IsDir:     record.IsDir(),
				Timestamp: 0,
			}
			continue
		}
		normalizedRecords[i] = FileRecord{
			Name:      record.Name(),
			IsDir:     record.IsDir(),
			Timestamp: info.ModTime().UnixNano(),
		}
	}

	if err := json.NewEncoder(w).Encode(&normalizedRecords); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UploadFileToDownloadsRepository(w http.ResponseWriter, r *http.Request) {
	_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_FILES_UPLOAD_TO_DOWNLOADS_REPO, nil)
	if reject {
		return
	}

	fileName := r.PathValue("fileName")

	file, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := os.WriteFile(fileName, file, 0666); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DownloadFileFromUploadsRepository(w http.ResponseWriter, r *http.Request) {
	// Resolve the absolute path of the uploads directory
	uploadsDir, err := filepath.Abs(*config.UploadsDirectoryPath)
	if err != nil {
		log.Printf("Error resolving uploads directory: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Extract and validate the file name from the request
	fileName := r.PathValue("fileName")
	if fileName == "" {
		http.Error(w, "Missing file name", http.StatusBadRequest)
		return
	}

	// Construct and clean the full file path
	fullPath := filepath.Join(uploadsDir, fileName)
	cleanPath := filepath.Clean(fullPath)

	// Prevent directory traversal by ensuring the path stays within uploadsDir
	if !strings.HasPrefix(cleanPath, uploadsDir+string(filepath.Separator)) {
		http.Error(w, "Invalid file path", http.StatusForbidden)
		return
	}

	// Check if the path exists and is a file
	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		log.Printf("Error stating file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if info.IsDir() {
		// TODO: Log and stuff.
		http.Error(w, "Directories cannot be downloaded", http.StatusForbidden)
		return
	}

	// Open the file for streaming
	file, err := os.Open(cleanPath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set headers for file download
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	if contentType := mime.TypeByExtension(filepath.Ext(fileName)); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	// Stream the file to the client
	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("Error streaming file: %v", err)
	}
}
