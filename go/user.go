package main

import (
	"database/sql"
	"time"
)

var userStmt *sql.Stmt

type timestamp time.Time

func (t timestamp) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(`"2006-01-02T15:04:05Z"`)), nil
}

type User struct {
	Id        string    `json:"id"`
	Email     string    `json:"email,omitempty"`
	CreatedAt timestamp `json:"created_at"`
	UpdatedAt timestamp `json:"updated_at"`
	Name      string    `json:"name,omitempty"`
	Admin     bool      `json:"admin"`
	Active    bool      `json:"-"`
}

func userInit(db *sql.DB) {
	var err error
	userStmt, err = db.Prepare("select id, email, created_at, updated_at, name, admin, active from users where apikey = $1")
	if err != nil {
		panic(err)
	}
}

func GetUserByApiKey(db *sql.DB, key Credentials) (*User, error) {
	rows, err := userStmt.Query(string(key))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var (
			user                 User
			id, email, name      []byte
			createdAt, updatedAt time.Time
		)
		err = rows.Scan(&id, &email, &createdAt, &updatedAt, &name, &user.Admin, &user.Active)
		if err != nil {
			return nil, err
		}
		user.Id = string(id)
		user.Email = string(email)
		user.Name = string(name)
		user.CreatedAt = timestamp(createdAt)
		user.UpdatedAt = timestamp(updatedAt)
		return &user, nil
	}
	return nil, nil
}
