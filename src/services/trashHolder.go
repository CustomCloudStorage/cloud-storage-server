package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/utils"
)

type TrashService interface {
	PermanentDeleteFile(ctx context.Context, userID, fileID int) error
	PermanentDeleteFolder(ctx context.Context, userID, folderID int) error
}

type trashService struct {
	trashRepository repositories.TrashRepository
	storageDir      string
}

func NewTrashService(trashRepo repositories.TrashRepository, cfg ServiceConfig) TrashService {
	svc := &trashService{
		trashRepository: trashRepo,
		storageDir:      cfg.StorageDir,
	}
	go svc.purgeLoop()
	return svc
}

func (s *trashService) PermanentDeleteFile(
	ctx context.Context, userID, fileID int) error {
	phys, err := s.trashRepository.PermanentDeleteFile(ctx, userID, fileID)
	if err != nil {
		return err
	}
	path := filepath.Join(s.storageDir, phys)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return utils.DetermineFSError(err, fmt.Sprintf("remove file %s", phys))
	}
	return nil
}

func (s *trashService) PermanentDeleteFolder(ctx context.Context, userID, folderID int) error {
	physList, err := s.trashRepository.PermanentDeleteFolder(ctx, userID, folderID)
	if err != nil {
		return err
	}

	for _, phys := range physList {
		path := filepath.Join(s.storageDir, phys)
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return utils.DetermineFSError(err, fmt.Sprintf("remove file %s", phys))
		}
	}
	return nil
}

func (s *trashService) purgeLoop() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		s.purgeOnce()
	}
}

func (s *trashService) purgeOnce() {
	cutoff := time.Now().Add(-30 * 24 * time.Hour)

	files, err := s.trashRepository.ListFilesToPurge(context.Background(), cutoff)
	if err != nil {
		fmt.Printf("trash GC: failed to list files to purge: %v\n", err)
	} else {
		for _, f := range files {
			path := filepath.Join(s.storageDir, f.PhysicalName)
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				fmt.Printf("trash GC: failed to remove file %s: %v\n", path, err)
			}
			if err := s.trashRepository.HardDeleteFileByID(context.Background(), f.ID); err != nil {
				fmt.Printf("trash GC: failed to hard delete file record %d: %v\n", f.ID, err)
			}
		}
	}

	folders, err := s.trashRepository.ListFoldersToPurge(context.Background(), cutoff)
	if err != nil {
		fmt.Printf("trash GC: failed to list folders to purge: %v\n", err)
	} else {
		for _, fld := range folders {
			if err := s.trashRepository.HardDeleteFolderByID(context.Background(), fld.ID); err != nil {
				fmt.Printf("trash GC: failed to hard delete folder record %d: %v\n", fld.ID, err)
			}
		}
	}
}
