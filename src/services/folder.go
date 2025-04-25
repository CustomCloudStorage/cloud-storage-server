package services

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (s *folderService) DownloadFolder(ctx context.Context, userID, folderID int) (io.ReadCloser, string, error) {
	folder, err := s.folderRepository.GetByID(ctx, folderID, userID)
	if err != nil {
		return nil, "", fmt.Errorf("folder not found or access denied: %w", err)
	}
	archiveName := folder.Name + ".zip"

	files, err := s.fileRepository.ListFilesRecursive(ctx, userID, folderID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list files: %w", err)
	}

	pr, pw := io.Pipe()
	go func() {
		zw := zip.NewWriter(pw)
		defer zw.Close()
		defer pw.Close()
		for _, f := range files {
			select {
			case <-ctx.Done():
				pw.CloseWithError(ctx.Err())
				return
			default:
			}
			path := filepath.Join(s.storageDir, f.PhysicalName)
			in, err := os.Open(path)
			if err != nil {
				pw.CloseWithError(fmt.Errorf("open %s: %w", path, err))
				return
			}
			info, _ := in.Stat()
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				in.Close()
				pw.CloseWithError(err)
				return
			}
			header.Name = f.RelativePath
			header.Method = zip.Deflate
			w, err := zw.CreateHeader(header)
			if err != nil {
				in.Close()
				pw.CloseWithError(err)
				return
			}
			if _, err := io.Copy(w, in); err != nil {
				in.Close()
				pw.CloseWithError(err)
				return
			}
			in.Close()
		}
	}()
	return pr, archiveName, nil
}
