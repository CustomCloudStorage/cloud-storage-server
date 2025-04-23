package service

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/CustomCloudStorage/types"
)

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
