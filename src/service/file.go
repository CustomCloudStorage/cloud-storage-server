package service

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/CustomCloudStorage/types"
	"github.com/google/uuid"
)

func (s *file) UploadFile(ctx context.Context, userID int, folderID *int, fileName string, fileSize int64, fileData io.Reader) error {
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

	return nil
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

	return s.repository.File.Delete(ctx, id, userID)
}
