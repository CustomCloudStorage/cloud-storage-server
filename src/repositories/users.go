package repositories

import (
	"context"

	"github.com/CustomCloudStorage/types"
	"github.com/CustomCloudStorage/utils"
)

const (
	get_user      = "SELECT id, name, email, role, storage-limit, last-update FROM users"
	get_all_users = "SELECT id, name, email, role, storage-limit FROM users"
)

func (postgres *Postgres) GetUser(ctx context.Context, id string) (*types.User, error) {
	var user types.User

	err := postgres.Db.QueryRowContext(ctx, get_user).Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.StorageLimit, &user.LastUpdate)
	if err != nil {
		return nil, utils.DetermineSQLError(err, "get data by id")
	}

	return &user, nil
}

func (postgres *Postgres) GetAllUsers(ctx context.Context) ([]types.User, error) {
	var users []types.User

	rows, err := postgres.Db.QueryContext(ctx, get_all_users)
	if err != nil {
		return nil, utils.DetermineSQLError(err, "get all data")
	}
	defer rows.Close()

	for rows.Next() {
		var user types.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.StorageLimit); err != nil {
			return nil, utils.DetermineSQLError(err, "scan row")
		}
		users = append(users, user)
	}

	return users, nil
}
