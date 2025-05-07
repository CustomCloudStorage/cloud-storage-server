package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
	"github.com/google/uuid"
)

func (s *uploadService) InitSession(ctx context.Context, session *types.UploadSession) error {
	if err := s.userRepository.ReserveStorage(ctx, session.UserID, session.TotalSize); err != nil {
		return err
	}

	id := uuid.New()
	session.ID = id
	if err := s.uploadSessionRepository.Create(ctx, session); err != nil {
		return err
	}

	path := filepath.Join(s.temp, id.String())
	if err := os.MkdirAll(path, 0o755); err != nil {
		return utils.DetermineFSError(err, fmt.Sprintf("mkdir temp dir %s", path))
	}
	return nil
}

func (s *uploadService) UploadPart(ctx context.Context, sessionID uuid.UUID, partNumber int, data io.Reader) error {
	session, err := s.uploadSessionRepository.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if partNumber < 1 || partNumber > session.TotalParts {
		return utils.ErrBadRequest.New("invalid part number %d", partNumber)
	}

	partPath := filepath.Join(s.temp, sessionID.String(), fmt.Sprintf("%05d.part", partNumber))
	f, err := os.Create(partPath)
	if err != nil {
		return utils.DetermineFSError(err, fmt.Sprintf("create part file %s", partPath))
	}
	defer f.Close()

	n, err := io.Copy(f, data)
	if err != nil {
		return utils.ErrInternal.Wrap(err, "write part %d data", partNumber)
	}

	part := &types.UploadPart{
		SessionID:  sessionID,
		PartNumber: partNumber,
		Size:       n,
	}
	if err := s.uploadPartRepository.Create(ctx, part); err != nil {
		return err
	}
	return nil
}

func (s *uploadService) GetProgress(ctx context.Context, sessionID uuid.UUID) (int64, int, error) {
	session, err := s.uploadSessionRepository.GetByID(ctx, sessionID)
	if err != nil {
		return 0, 0, err
	}
	parts, err := s.uploadPartRepository.ListBySession(ctx, sessionID)
	if err != nil {
		return 0, 0, err
	}
	var total int64
	for _, p := range parts {
		total += p.Size
	}
	return total, session.TotalParts, nil
}

func (s *uploadService) Complete(ctx context.Context, sessionID uuid.UUID) (*types.File, error) {
	session, err := s.uploadSessionRepository.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	physical := uuid.New().String()
	finalPath := filepath.Join(s.storageDir, physical)
	out, err := os.Create(finalPath)
	if err != nil {
		return nil, utils.DetermineFSError(err, fmt.Sprintf("create final file %s", finalPath))
	}
	defer out.Close()

	for i := 1; i <= session.TotalParts; i++ {
		partPath := filepath.Join(s.temp, sessionID.String(), fmt.Sprintf("%05d.part", i))
		in, err := os.Open(partPath)
		if err != nil {
			return nil, utils.DetermineFSError(err, fmt.Sprintf("open part %d", i))
		}
		if _, err := io.Copy(out, in); err != nil {
			in.Close()
			return nil, utils.ErrInternal.Wrap(err, "copy part %d", i)
		}
		in.Close()
	}

	fileMeta := &types.File{
		UserID:       session.UserID,
		FolderID:     session.FolderID,
		Name:         session.Name,
		Extension:    session.Extension,
		Size:         session.TotalSize,
		PhysicalName: physical,
	}
	if err := s.fileRepository.Create(ctx, fileMeta); err != nil {
		if err := s.userRepository.ReleaseStorage(ctx, session.UserID, session.TotalSize); err != nil {
			return nil, err
		}
		return nil, err
	}

	user, err := s.userRepository.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}
	if err := s.userRepository.UpdateUsedStorage(ctx, session.UserID, user.Account.UsedStorage+fileMeta.Size); err != nil {
		return nil, err
	}

	tempDir := filepath.Join(s.temp, sessionID.String())
	if err := os.RemoveAll(tempDir); err != nil {
		return nil, utils.DetermineFSError(err, fmt.Sprintf("remove temp dir %s", tempDir))
	}
	if err := s.uploadPartRepository.DeleteBySession(ctx, sessionID); err != nil {
		return nil, err
	}
	if err := s.uploadSessionRepository.Delete(ctx, sessionID); err != nil {
		return nil, err
	}

	return fileMeta, nil
}

func (s *uploadService) Abort(ctx context.Context, sessionID uuid.UUID) error {
	session, err := s.uploadSessionRepository.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if err := s.userRepository.ReleaseStorage(ctx, session.UserID, session.TotalSize); err != nil {
		return err
	}

	tempDir := filepath.Join(s.temp, sessionID.String())
	if err := os.RemoveAll(tempDir); err != nil {
		return utils.DetermineFSError(err, fmt.Sprintf("remove temp dir %s", tempDir))
	}
	if err := s.uploadPartRepository.DeleteBySession(ctx, sessionID); err != nil {
		return err
	}
	if err := s.uploadSessionRepository.Delete(ctx, sessionID); err != nil {
		return err
	}
	return nil
}

func (s *uploadService) purgeLoop() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		s.purgeOnce()
	}
}

func (s *uploadService) purgeOnce() {
	sessions, err := s.uploadSessionRepository.ListOlderThan(context.Background(), 7*24*time.Hour)
	if err != nil {
		fmt.Printf("upload GC: cannot list stale sessions: %v\n", err)
		return
	}

	for _, sess := range sessions {
		dir := filepath.Join(s.temp, sess.ID.String())
		if err := os.RemoveAll(dir); err != nil {
			fmt.Printf("upload GC: remove tmp dir %s: %v\n", dir, err)
		}
		if err := s.uploadSessionRepository.Delete(context.Background(), sess.ID); err != nil {
			fmt.Printf("upload GC: delete session %s: %v\n", sess.ID, err)
		}
	}
}
