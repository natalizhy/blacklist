package models

type User struct {
	ID        int64  `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Country   string `db:"country"`
	Phone     string `db:"phone"`
	Info      string `db:"info"`
}