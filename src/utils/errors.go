package utils

import (
	"database/sql"

	"github.com/joomcode/errorx"
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
	ErrGet        = ErrHandler.NewType("get_error")
	ErrJsonEncode = ErrHandler.NewType("json_encode_error")

	ErrRepository = errorx.NewNamespace("repository")
	ErrNotFound   = ErrRepository.NewType("not_found_error")
	ErrSql        = ErrRepository.NewType("sql_error")
)

func DetermineSQLError(err error, context string) error {
	if err == sql.ErrNoRows {
		return ErrNotFound.New("data not found: %s", context)
	}
	return ErrSql.Wrap(err, context)
}
