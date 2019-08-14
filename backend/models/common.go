package models

type User struct {
	ID        int64  `db:"id"`
	FirstName string `db:"first_name" validate:"required,alpha"`
	LastName  string `db:"last_name" validate:"required"`
	Country   string `db:"country" validate:"required"`
	Phone     string `db:"phone" validate:"required"`
	Info      string `db:"info" validate:"required"`
	Photo     string `db:"photo" validate:"required"`
}
