package users

import (
  "time"
  "go-teach-me/database"
)

type User struct {
  UserID int
  FirstName string
  LastName string
  Email string
  CreatedAt time.Time
  UpdatedAt time.Time
}

func GetUser(email string, password string) (User) {
  var user_id int
  var first_name string
  var last_name string
  var created_at time.Time
  var updated_at time.Time
  database.DBCon.QueryRow("SELECT user_id, first_name, last_name, created_at, updated_at FROM teachme.users WHERE email = $1 AND password = $2", email, password).Scan(&user_id, &first_name, &last_name, &created_at, &updated_at)
  return User{UserID: user_id, FirstName: first_name, LastName: last_name, Email: email, CreatedAt: created_at, UpdatedAt: updated_at}
}

func InsertUser(firstName string, lastName string, email string, password string, userGroupID int) (error){
  stmt, err := database.DBCon.Prepare("INSERT INTO teachme.users(first_name, last_name, email, password, user_group_id) VALUES($1, $2, $3, $4, $5)")
  if err != nil {
    return err
  }
  _, err = stmt.Exec(firstName, lastName, email, password, userGroupID)
  if err != nil {
    return err
  }
  return nil
}
