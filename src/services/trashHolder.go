package services

import (
	"context"
	"os"
	"path/filepath"
	"time"
)

func (s *trashService) PermanentDeleteFile(
	ctx context.Context, userID, fileID int) error {
	phys, err := s.trashRepository.PermanentDeleteFile(ctx, userID, fileID)
	if err != nil {
		return err
	}
	os.Remove(filepath.Join(s.storageDir, phys))
	return nil
}

func (s *trashService) PermanentDeleteFolder(ctx context.Context, userID, folderID int) error {
	physList, err := s.trashRepository.PermanentDeleteFolder(ctx, userID, folderID)
	if err != nil {
		return err
	}

	for _, name := range physList {
		os.Remove(filepath.Join(s.storageDir, name))
	}
	return nil
}

func (s *trashService) purgeLoop() {
	ticker := time.NewTicker(s.cfg.CleanupInterval)
	defer ticker.Stop()
	for range ticker.C {
		s.purgeOnce()
	}
}

func (s *trashService) purgeOnce() {
	cutoff := time.Now().Add(-s.cfg.TTL)

	files, _ := s.trashRepository.ListFilesToPurge(context.Background(), cutoff)
	for _, f := range files {
		os.Remove(filepath.Join(s.storageDir, f.PhysicalName))
		s.trashRepository.HardDeleteFileByID(context.Background(), f.ID)
	}

	folders, _ := s.trashRepository.ListFoldersToPurge(context.Background(), cutoff)
	for _, fld := range folders {
		s.trashRepository.HardDeleteFolderByID(context.Background(), fld.ID)
	}
}
