package services

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/CustomCloudStorage/utils"
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
				pw.CloseWithError(utils.DetermineFSError(err, "open "+path))
				return
			}
			info, err := in.Stat()
			if err != nil {
				in.Close()
				pw.CloseWithError(utils.DetermineFSError(err, "stat "+path))
				return
			}
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				in.Close()
				pw.CloseWithError(utils.ErrInternal.Wrap(err, "create zip header for %s", path))
				return
			}
			header.Name = f.RelativePath
			header.Method = zip.Deflate
			zwEntry, err := zw.CreateHeader(header)
			if err != nil {
				in.Close()
				pw.CloseWithError(utils.ErrInternal.Wrap(err, "create zip entry %s", f.RelativePath))
				return
			}
			if _, err := io.Copy(zwEntry, in); err != nil {
				in.Close()
				pw.CloseWithError(utils.ErrInternal.Wrap(err, "write data for %s", path))
				return
			}
			in.Close()
		}
	}()
	return pr, archiveName, nil
}
