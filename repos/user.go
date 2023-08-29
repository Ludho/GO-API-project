package repos

import (
	"context"
	"fmt"
	"trades/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(context.Context, *models.User) error
	GetUserByUsername(context.Context, string) (*models.User, error)
}

type userRepositoryImpl struct {
	connection *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) UserRepository {
	return &userRepositoryImpl{
		connection: conn,
	}
}

const SQL_INSERT_USER = `
	insert into public.user 
	values( $1, $2, $3)
`

func (u *userRepositoryImpl) CreateUser(c context.Context, user *models.User) error {

	_, err := u.connection.Query(c, SQL_INSERT_USER, user.Id, user.Username, user.Password)

	if err != nil {
		return fmt.Errorf("error during query to create user: %v", err)
	}
	return nil
}

const SQL_GET_USER = `
						select 
							u.id,
							u.username,
							u.pass
						from
							public.user as u
						where u.username = $1;
`

func (u *userRepositoryImpl) GetUserByUsername(c context.Context, username string) (*models.User, error) {

	rows, err := u.connection.Query(c, SQL_GET_USER, username)

	if err != nil {
		return nil, fmt.Errorf("error during query to get user: %v", err)
	}

	if rows.Next() {
		user := &models.User{}
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Password,
		)

		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, fmt.Errorf(`no user with username "%s" found`, username)

}
