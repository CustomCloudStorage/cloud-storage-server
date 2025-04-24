package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/CustomCloudStorage/types"
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

	path := filepath.Join(s.tmpUpload, id.String())
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("mkdir tmp session dir: %w", err)
	}
	return nil
}

func (s *uploadService) UploadPart(ctx context.Context, sessionID uuid.UUID, partNumber int, data io.Reader) error {
	session, err := s.uploadSessionRepository.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("get session: %w", err)
	}
	if partNumber < 1 || partNumber > session.TotalParts {
		return fmt.Errorf("invalid part number %d", partNumber)
	}

	partPath := filepath.Join(s.tmpUpload, sessionID.String(), fmt.Sprintf("%05d.part", partNumber))
	f, err := os.Create(partPath)
	if err != nil {
		return fmt.Errorf("create part file: %w", err)
	}
	defer f.Close()

	n, err := io.Copy(f, data)
	if err != nil {
		return fmt.Errorf("write part data: %w", err)
	}

	part := &types.UploadPart{
		SessionID:  sessionID,
		PartNumber: partNumber,
		Size:       n,
	}
	if err := s.uploadPartRepository.Create(ctx, part); err != nil {
		return fmt.Errorf("save part metadata: %w", err)
	}
	return nil
}

func (s *uploadService) GetProgress(ctx context.Context, sessionID uuid.UUID) (int64, int, error) {
	session, err := s.uploadSessionRepository.GetByID(ctx, sessionID)
	if err != nil {
		return 0, 0, fmt.Errorf("get session: %w", err)
	}
	parts, err := s.uploadPartRepository.ListBySession(ctx, sessionID)
	if err != nil {
		return 0, 0, fmt.Errorf("list parts: %w", err)
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
		return nil, fmt.Errorf("get session: %w", err)
	}

	physical := uuid.New().String()
	finalPath := filepath.Join(s.storageDir, physical)
	out, err := os.Create(finalPath)
	if err != nil {
		return nil, fmt.Errorf("create final file: %w", err)
	}
	defer out.Close()

	for i := 1; i <= session.TotalParts; i++ {
		partPath := filepath.Join(s.tmpUpload, sessionID.String(), fmt.Sprintf("%05d.part", i))
		in, err := os.Open(partPath)
		if err != nil {
			return nil, fmt.Errorf("open part %d: %w", i, err)
		}
		if _, err := io.Copy(out, in); err != nil {
			in.Close()
			return nil, fmt.Errorf("copy part %d: %w", i, err)
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
		return nil, fmt.Errorf("save file metadata: %w", err)
	}

	user, err := s.userRepository.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if err := s.userRepository.UpdateUsedStorage(ctx, session.UserID, user.Account.UsedStorage+fileMeta.Size); err != nil {
		return nil, fmt.Errorf("update used_storage: %w", err)
	}

	if err := os.RemoveAll(filepath.Join(s.tmpUpload, sessionID.String())); err != nil {
		return nil, err
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

	os.RemoveAll(filepath.Join(s.tmpUpload, sessionID.String()))
	if err := s.uploadPartRepository.DeleteBySession(ctx, sessionID); err != nil {
		return fmt.Errorf("delete parts metadata: %w", err)
	}
	if err := s.uploadSessionRepository.Delete(ctx, sessionID); err != nil {
		return fmt.Errorf("delete session metadata: %w", err)
	}
	return nil
}
