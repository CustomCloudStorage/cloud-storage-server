package utils

import (
	"errors"
	"os"
	"strings"
	"syscall"

	"github.com/joomcode/errorx"
	"gorm.io/gorm"
)

var (
	Namespace = errorx.NewNamespace("app_error")

	ErrBadRequest = Namespace.NewType("bad_request")
	ErrNotFound   = Namespace.NewType("not_found")
	ErrConflict   = Namespace.NewType("conflict")
	ErrInternal   = Namespace.NewType("internal")
)

func DetermineSQLError(err error, context string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound.Wrap(err, "data not found: %s", context)
	}
	if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return ErrConflict.Wrap(err, "data conflict: %s", context)
	}
	return ErrInternal.Wrap(err, "sql error: %s", context)
}

func DetermineFSError(err error, context string) error {
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return ErrNotFound.Wrap(err, "file not found: %s", context)
	}
	if os.IsPermission(err) {
		return ErrInternal.Wrap(err, "permission denied: %s", context)
	}
	var pe *os.PathError
	if errors.As(err, &pe) {
		if errno, ok := pe.Err.(syscall.Errno); ok {
			switch errno {
			case syscall.ENOTEMPTY:
				return ErrConflict.Wrap(err, "directory not empty: %s", context)
			case syscall.EMFILE:
				return ErrInternal.Wrap(err, "too many open files: %s", context)
			}
		}
	}
	return ErrInternal.Wrap(err, "I/O error: %s", context)
}
