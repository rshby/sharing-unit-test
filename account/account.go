package account

import "database/sql"

type Account struct {
	Id       int            `json:"id,omitempty"`
	Email    string         `json:"email,omitempty"`
	Username string         `json:"username,omitempty"`
	Password string         `json:"password,omitempty"`
	FullName sql.NullString `json:"full_name,omitempty"`
	Gender   sql.NullString `json:"gender,omitempty"`
}
