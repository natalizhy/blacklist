package repositories

import (
	//"database/sql"
	"fmt"
	"github.com/natalizhy/blacklist/backend/models"

	//"github.com/jmoiron/sqlx"
)

func DBUser(user models.User) {

	res, err := DB.Exec("INSERT INTO `customers` (`first_name`, `last_name`, `country`, `phone`, `info`) VALUES (?, ?, ?, ?, ?)", user.FirstName, user.LastName, user.Country, user.Phone, user.Info)

	if err != nil {
		panic(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Created user with id:", id)
}

//func DBUpdate() {
//
//	res, err := database.Exec("INSERT INTO `customers` (`first_name`, `last_name`, `country`, `phone`, `info`, `photo`) VALUES (?, ?, ?, ?, ?, ?, ?)", `FirstName`, `LastName`, `Country`, `Phone`, `Info`, `Photo`)
//
//	if err != nil {
//		panic(err)
//	}
//	id, err := res.LastInsertId()
//	if err != nil {
//		panic(err)
//	}
//
//	_, err = database.Exec("UPDATE customers set name=\"John\" where id=?", id)
//	if err != nil {
//		panic(err)
//	}
//}
//func DBDelete() {
//	res, err := database.Exec("INSERT INTO `customers` (`first_name`, `last_name`, `country`, `phone`, `info`, `photo`) VALUES (?, ?, ?, ?, ?, ?, ?)", `FirstName`, `LastName`, `Country`, `Phone`, `Info`, `Photo`)
//
//	if err != nil {
//		panic(err)
//	}
//	id, err := res.LastInsertId()
//	if err != nil {
//		panic(err)
//	}
//
//	_, err = database.Exec("DELETE FROM customers where id=?", id)
//	if err != nil {
//		panic(err)
//	}
//}
