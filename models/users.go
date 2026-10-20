package models

import (
	"errors"
	"fmt"

	"github.com/zahra-ove/go-restapi-sample/db"
	"github.com/zahra-ove/go-restapi-sample/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (user *User) Save() error {
	query := `INSERT INTO users(email, password) VALUES (?, ?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()
	hashedPassword, err := utils.HashedPassword(user.Password)
	if err != nil {
		return err
	}

	userData, err := stmt.Exec(user.Email, hashedPassword)
	user.ID, err = userData.LastInsertId()

	if err != nil {
		return err
	}

	return err
}

func (u *User) ValidateCredentials() error {
	query := `SELECT id, password FROM users WHERE email = ?`
	row := db.DB.QueryRow(query, u.Email)
	fmt.Print("here0 email:\n", u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)

	if err != nil {
		return errors.New("Invalid Credentials")
	}

	passwordIsValid := utils.CheckHashedPassword(retrievedPassword, u.Password)

	if !passwordIsValid {
		return errors.New("Invalid Credentials")
	}

	return nil
}
