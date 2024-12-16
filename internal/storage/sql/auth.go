package sql

import (
	"auth_rest/internal/utils"
	_ "embed"
	"github.com/jackc/pgx/v5"
)

type AuthRepository interface {
	GetRefreshToken(guid string) (string, error)
	PostRefreshToken(guid string, refreshToken string) error
	DeleteRefreshToken(guid string) error
}

type AuthPostgres struct {
	db  *pgx.Conn
	ctx utils.AppContext
}

func NewAuthPostgres(db *pgx.Conn, ctx utils.AppContext) *AuthPostgres {
	return &AuthPostgres{
		db:  db,
		ctx: ctx,
	}
}

////go:embed query/GetUser.sql
//var getUser string
//
//func (r *AuthPostgres) GetUser(guid string) (models.User, error) {
//	var user models.User
//
//	row := r.db.QueryRow(r.ctx.Ctx, getUser, guid)
//	if err := row.Scan(&user); err != nil {
//		return models.User{}, err
//	}
//
//	return user, nil
//}

//go:embed query/PostRefreshToken.sql
var postRefreshToken string

func (r *AuthPostgres) PostRefreshToken(guid string, refreshToken string) error {
	_, err := r.db.Exec(r.ctx.Ctx, postRefreshToken, guid, refreshToken)
	if err != nil {
		return err
	}
	return nil
}

//go:embed query/GetRefreshToken.sql
var getRefreshToken string

func (r *AuthPostgres) GetRefreshToken(guid string) (string, error) {
	var refreshToken string
	row := r.db.QueryRow(r.ctx.Ctx, getRefreshToken, guid)
	if err := row.Scan(&refreshToken); err != nil {
		return "", err
	}

	return refreshToken, nil
}

//go:embed query/DeleteRefreshToken.sql
var deleteRefreshToken string

func (r *AuthPostgres) DeleteRefreshToken(guid string) error {
	_, err := r.db.Exec(r.ctx.Ctx, deleteRefreshToken, guid)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthPostgres) Close() error {
	err := r.db.Close(r.ctx.Ctx)
	if err != nil {
		return err
	}

	return nil
}
