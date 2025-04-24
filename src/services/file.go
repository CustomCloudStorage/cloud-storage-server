package services

import (
	"context"
	"os"
	"path/filepath"
)

func (s *fileService) DeleteFile(ctx context.Context, id int, userID int) error {
	file, err := s.fileRepository.GetByID(ctx, id, userID)
	if err != nil {
		return nil
	}

	filePath := filepath.Join(s.storageDir, file.PhysicalName)
	if err := os.Remove(filePath); err != nil {
		return err
	}

	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	newUsed := user.Account.UsedStorage - file.Size
	if newUsed < 0 {
		newUsed = 0
	}
	if err := s.userRepository.UpdateUsedStorage(ctx, userID, newUsed); err != nil {
		return err
	}

	return nil
}
