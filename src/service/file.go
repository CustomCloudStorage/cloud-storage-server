package service

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/CustomCloudStorage/types"
	"github.com/google/uuid"
)

func (s *file) UploadFile(ctx context.Context, userID int, folderID *int, fileName string, fileSize int64, fileData io.Reader) error {
	user, err := s.repository.User.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.Account.UsedStorage+fileSize > user.Account.StorageLimit {
		return nil
	}

	extension := filepath.Ext(fileName)
	baseName := strings.TrimSuffix(fileName, extension)

	physicalName := uuid.New().String()
	destPath := filepath.Join(s.storageDir, physicalName)

	outFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, fileData); err != nil {
		return err
	}

	newFile := &types.File{
		UserID:       userID,
		FolderID:     folderID,
		Name:         baseName,
		Extension:    extension,
		Size:         fileSize,
		PhysicalName: physicalName,
	}

	if err := s.repository.File.Create(ctx, newFile); err != nil {
		os.Remove(destPath)
		return err
	}

	if err := s.repository.User.UpdateUsedStorage(ctx, userID, user.Account.UsedStorage+fileSize); err != nil {
		return err
	}

	return nil
}

func (s *file) DownloadFile(ctx context.Context, userID int, fileID int) (*types.DownloadedFile, error) {
	fileMeta, err := s.repository.File.GetByID(ctx, fileID, userID)
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(s.storageDir, fileMeta.PhysicalName)

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 512)
	n, err := f.Read(buffer)
	if err != nil && err != io.EOF {
		f.Close()
		return nil, err
	}
	contentType := http.DetectContentType(buffer[:n])

	if _, err = f.Seek(0, io.SeekStart); err != nil {
		f.Close()
		return nil, err
	}

	downloadedFileName := fileMeta.Name + fileMeta.Extension

	return &types.DownloadedFile{
		Reader:      f,
		FileName:    downloadedFileName,
		ContentType: contentType,
		FileSize:    fileMeta.Size,
		ModTime:     fileMeta.UpdatedAt,
	}, nil
}

func (s *file) DeleteFile(ctx context.Context, id int, userID int) error {
	file, err := s.repository.File.GetByID(ctx, id, userID)
	if err != nil {
		return nil
	}

	filePath := filepath.Join(s.storageDir, file.PhysicalName)
	if err := os.Remove(filePath); err != nil {
		return err
	}

	if err := s.repository.File.Delete(ctx, id, userID); err != nil {
		return err
	}

	user, err := s.repository.User.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	newUsed := user.Account.UsedStorage - file.Size
	if newUsed < 0 {
		newUsed = 0
	}
	if err := s.repository.User.UpdateUsedStorage(ctx, userID, newUsed); err != nil {
		return err
	}

	return nil
}
