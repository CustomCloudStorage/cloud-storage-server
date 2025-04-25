package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/CustomCloudStorage/types"
)

const tokenParts = 4

func (s *fileService) GenerateDownloadURL(ctx context.Context, userID, fileID int) (string, error) {
	if _, err := s.fileRepository.GetByID(ctx, fileID, userID); err != nil {
		return "", err
	}

	expiry := time.Now().Add(s.urlTtlSeconds).Unix()
	payload := fmt.Sprintf("%d:%d:%d", userID, fileID, expiry)

	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write([]byte(payload))
	sig := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	token := base64.URLEncoding.EncodeToString([]byte(payload + ":" + sig))

	q := url.QueryEscape(token)
	return fmt.Sprintf("%s/files/download?token=%s", s.host, q), nil
}

func (s *fileService) ValidateDownloadToken(token string) (int, int, error) {
	raw, err := url.QueryUnescape(token)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid token encoding: %w", err)
	}
	data, err := base64.URLEncoding.DecodeString(raw)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid token: %w", err)
	}

	parts := strings.Split(string(data), ":")
	if len(parts) != tokenParts {
		return 0, 0, fmt.Errorf("malformed token")
	}

	uID, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid user in token")
	}
	fID, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid file in token")
	}
	expiry, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid expiry in token")
	}

	if time.Now().Unix() > expiry {
		return 0, 0, fmt.Errorf("token expired")
	}

	payload := strings.Join(parts[:3], ":")
	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write([]byte(payload))
	expectedSig := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expectedSig), []byte(parts[3])) {
		return 0, 0, fmt.Errorf("invalid signature")
	}
	return uID, fID, nil
}

func (s *fileService) DownloadFile(ctx context.Context, userID int, fileID int) (*types.DownloadedFile, error) {
	fileMeta, err := s.fileRepository.GetByID(ctx, fileID, userID)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(s.storageDir, fileMeta.PhysicalName)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(fileMeta.Extension)
	var ctype string
	switch ext {
	case ".jpg", ".jpeg":
		ctype = "image/jpeg"
	case ".png":
		ctype = "image/png"
	case ".gif":
		ctype = "image/gif"
	case ".mp4":
		ctype = "video/mp4"
	case ".webm":
		ctype = "video/webm"
	case ".pdf":
		ctype = "application/pdf"
	default:
		buf := make([]byte, 512)
		n, _ := f.Read(buf)
		ctype = http.DetectContentType(buf[:n])
		f.Seek(0, io.SeekStart)
	}

	return &types.DownloadedFile{
		Reader:      f,
		FileName:    fileMeta.Name + fileMeta.Extension,
		ContentType: ctype,
		FileSize:    fileMeta.Size,
		ModTime:     fileMeta.UpdatedAt,
	}, nil
}

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
