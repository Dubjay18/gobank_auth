package domain

import (
	"database/sql"
	"errors"
	"github.com/Dubjay18/gobank_auth/errs"
	"github.com/Dubjay18/gobank_auth/logger"
	"github.com/jmoiron/sqlx"
)

const (
	sqlVerifyQry = `SELECT username, u.customer_id, role, group_concat(a.account_id) as account_numbers FROM users u
                  LEFT JOIN accounts a ON a.customer_id = u.customer_id
                WHERE username = $1 and password = $2
                GROUP BY a.customer_id`
)

type AuthRepository interface {
	FindById(username string, password string) (*SLogin, *errs.AppError)
	GenerateAndSaveRefreshTokenToStore(token AuthToken) (string, *errs.AppError)
	RefreshTokenExists(refreshToken string) *errs.AppError
}

type AuthRepositoryDb struct {
	client *sqlx.DB
}

func (d AuthRepositoryDb) RefreshTokenExists(refreshToken string) *errs.AppError {
	sqlSelect := "select refresh_token from refresh_token_store where refresh_token = $1"
	var token string
	err := d.client.Get(&token, sqlSelect, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewAuthenticationError("refresh token not registered in the store")
		} else {
			logger.Error("Unexpected database error: " + err.Error())
			return errs.NewUnexpectedError("unexpected database error")
		}
	}
	return nil
}

func (d AuthRepositoryDb) GenerateAndSaveRefreshTokenToStore(authToken AuthToken) (string, *errs.AppError) {
	// generate the refresh token
	var appErr *errs.AppError
	var refreshToken string
	if refreshToken, appErr = authToken.newRefreshToken(); appErr != nil {
		return "", appErr
	}

	// store it in the store
	sqlInsert := "insert into refresh_token_store (refresh_token) values ($1)"
	_, err := d.client.Exec(sqlInsert, refreshToken)
	if err != nil {
		logger.Error("unexpected database error: " + err.Error())
		return "", errs.NewUnexpectedError("unexpected database error")
	}
	return refreshToken, nil
}
func (d AuthRepositoryDb) FindById(username, password string) (*SLogin, *errs.AppError) {
	var login SLogin

	err := d.client.Get(&login, sqlVerifyQry, username, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewAuthenticationError("invalid credentials")
		} else {
			logger.Error("Error while verifying login request from database: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &login, nil
}
func NewAuthRepository(client *sqlx.DB) AuthRepositoryDb {
	return AuthRepositoryDb{client}
}
