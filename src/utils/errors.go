package utils

import (
	"database/sql"

	"github.com/joomcode/errorx"
	"github.com/lib/pq"
)

var (
	ErrDatabase     = errorx.NewNamespace("database")
	ErrConnection   = ErrDatabase.NewType("connection_error")
	ErrMigration    = ErrDatabase.NewType("migration_error")
	ErrPingFailed   = ErrDatabase.NewType("ping_failed")
	ErrDriverCreate = ErrDatabase.NewType("driver_create_error")

	ErrConfig     = errorx.NewNamespace("config")
	ErrRead       = ErrConfig.NewType("read_error")
	ErrUnmarshal  = ErrConfig.NewType("unmarshal_error")
	ErrValidation = ErrConfig.NewType("validation_error")

	ErrHandler    = errorx.NewNamespace("handler")
	ErrJsonDecode = ErrHandler.NewType("json_decode_error")
	ErrJsonEncode = ErrHandler.NewType("json_encode_error")

	ErrRepository   = errorx.NewNamespace("repository")
	ErrNotFound     = ErrRepository.NewType("not_found_error")
	ErrSql          = ErrRepository.NewType("sql_error")
	ErrAlreadyExist = ErrRepository.NewType("already_exist_error")

	ErrDateTime = errorx.NewNamespace("date/time")
	ErrLocation = ErrDateTime.NewType("location_error")
	ErrFormat   = ErrDateTime.NewType("format_error")
)

func DetermineSQLError(err error, context string) error {
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "23505" {
			return ErrAlreadyExist.Wrap(err, "unique data is already exists: %s", context)
		}
	}
	if err == sql.ErrNoRows {
		return ErrNotFound.Wrap(err, "data not found: %s", context)
	}
	return ErrSql.Wrap(err, context)
}
