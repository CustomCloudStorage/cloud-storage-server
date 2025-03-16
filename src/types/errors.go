package types

import "github.com/joomcode/errorx"

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
)
