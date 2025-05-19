package services

import (
	"context"
	"fmt"
	"time"

	"github.com/CustomCloudStorage/repositories"
)

type TrashService interface {
	PermanentDeleteFile(ctx context.Context, userID, fileID int) error
	PermanentDeleteFolder(ctx context.Context, userID, folderID int) error
}

type trashService struct {
	trashRepository repositories.TrashRepository
	fileService     FileService
}

func NewTrashService(trashRepo repositories.TrashRepository, fileService FileService) TrashService {
	svc := &trashService{
		trashRepository: trashRepo,
		fileService:     fileService,
	}
	go svc.purgeLoop()
	return svc
}

func (s *trashService) PermanentDeleteFile(ctx context.Context, userID, fileID int) error {
	if err := s.trashRepository.PermanentDeleteFile(ctx, userID, fileID); err != nil {
		return err
	}
	if err := s.fileService.DeleteFile(ctx, fileID, userID); err != nil {
		return err
	}
	return nil
}

func (s *trashService) PermanentDeleteFolder(ctx context.Context, userID, folderID int) error {
	filesID, err := s.trashRepository.PermanentDeleteFolder(ctx, userID, folderID)
	if err != nil {
		return err
	}

	for _, id := range filesID {
		if err := s.fileService.DeleteFile(ctx, id, userID); err != nil {
			return err
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
			if err := s.fileService.DeleteFile(context.Background(), f.ID, f.UserID); err != nil {
				fmt.Printf("trash GC: failed to remove file: %v\n", err)
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
