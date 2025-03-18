package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
)

const (
	get_user = "SELECT id, name, email, role, storage-limit FROM users"
)

func (postgres *Postgres) GetUser(ctx context.Context, id string) (*types.User, error) {
	var user types.User

	err := postgres.Db.QueryRowContext(ctx, get_user).Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.StorageLimit)
	if err != nil {
		return nil, utils.DetermineSQLError(err, "get data by id")
	}

	return &user, nil
}
