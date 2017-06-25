package users

import (
	"errors"
	"fmt"
	"time"

	"go-teach-me/database"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID      int
	FirstName   string
	LastName    string
	Email       string
	UserGroupID int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func GetUser(email string, password string) (User, error) {
	var userID int
	var firstName string
	var lastName string
	var userGroupID int
	var createdAt time.Time
	var updatedAt time.Time

	var hashedPassword string
	resultRow := database.DBCon.QueryRow("SELECT user_id, first_name, last_name, password, user_group_id, created_at, updated_at FROM teachme.users WHERE email = $1", email)
  resultRow.Scan(&userID, &firstName, &lastName, &hashedPassword, &userGroupID, &createdAt, &updatedAt)
	if userID == 0 {
		return User{}, errors.New("User with provided email and password not found")
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return User{}, err
	}
	return User{UserID: userID, FirstName: firstName, LastName: lastName, Email: email, UserGroupID: userGroupID, CreatedAt: createdAt, UpdatedAt: updatedAt}, nil
}

func InsertUser(firstName string, lastName string, email string, password string, userGroupID int) error {
	hashedPass, err := hashPassword(password)
	fmt.Println(hashedPass)

	if err != nil {
		return err
	}
	stmt, err := database.DBCon.Prepare("INSERT INTO teachme.users(first_name, last_name, email, password, user_group_id) VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(firstName, lastName, email, hashedPass, userGroupID)
	if err != nil {
		return err
	}
	return nil
}
