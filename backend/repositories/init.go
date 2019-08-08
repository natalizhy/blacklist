package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitDB() {
	conn, err := sqlx.Connect("mysql", "root:@/users")
	if err != nil {

		panic(err)

	} else {
		fmt.Println("DB OK")
	}

	DB = conn
}