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

func (user User) Save() error {
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

	_, err = stmt.Exec(user.Email, hashedPassword)
	return err
}

func (u *User) ValidateCredentials() error {
	query := `SELECT id, password FROM users WHERE email = ?`
	row := db.DB.QueryRow(query, u.Email)
	fmt.Print("here0 email:\n", u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)

	if err != nil {
		fmt.Print("here1\n")
		fmt.Print("err:", err)
		return errors.New("Invalid Credentials")
	}

	passwordIsValid := utils.CheckHashedPassword(retrievedPassword, u.Password)

	if !passwordIsValid {
		fmt.Print("here2\n")
		fmt.Print("PasswordIsInvalid:", passwordIsValid, "\n")
		fmt.Print("hashed password:", u.Password, "\n")
		fmt.Print("password:", retrievedPassword, "\n")
		return errors.New("Invalid Credentials")
	}

	return nil
}
