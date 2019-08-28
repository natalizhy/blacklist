package models

type User struct {
	ID        int64  `db:"id"`
	FirstName string `db:"first_name" validate:"required,cyr,max=50,min=2"`
	LastName  string `db:"last_name" validate:"required,alpha,cyr,max=50,min=4"`
	CityID    int64  `db:"city_id"`
	Phone     string `db:"phone" validate:"required,numeric,max=50,min=4"`
	Info      string `db:"info" validate:"required,max=255,min=5"`
	Photo     string `db:"photo"`
}

type Signup struct {
	Login    string `db:"login" validate:"required,cyr,max=50,min=2"`
	Password string `db:"password"`
}
