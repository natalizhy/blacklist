package repositories

import (
	//"database/sql"
	"fmt"
	"github.com/natalizhy/blacklist/backend/models"

	//"github.com/jmoiron/sqlx"
)

func GetUsers() (users []models.User, err error) {
	rows, err := DB.Queryx("SELECT `id`, `first_name`, `last_name`, `country`, `phone`, `info` FROM `customers` ORDER BY id DESC")

	if err != nil {
		return
	}

	for rows.Next() {
		user := models.User{}
		_ = rows.StructScan(&user)

		users = append(users, user)
	}

	return
}

func GetUserById(userID int64) (user models.User, err error) {
	result := DB.QueryRowx("SELECT `id`, `first_name`, `last_name`, `country`, `phone`, `info` FROM `customers` WHERE id=?", userID)

	err = result.StructScan(&user)

	return
}

func UpdateUser(user models.User, userID int64) (err error) {
	_, err = DB.Exec("UPDATE `customers` SET `first_name`=?, `last_name`=?, `country`=?, `info`=?, `phone`=? WHERE id=?", user.FirstName, user.LastName, user.Country, user.Info, user.Phone, userID)
	return
}

func AddUser(user models.User) (newUserId int64, err error) {
	newUserId = 0

	res, err := DB.Exec("INSERT INTO `customers` (`first_name`, `last_name`, `country`, `phone`, `info`) VALUES (?, ?, ?, ?, ?)", user.FirstName, user.LastName, user.Country, user.Phone, user.Info)
	if err != nil {
		return
	}
	newUserId, err = res.LastInsertId()
	if err != nil {
		return
	}

	fmt.Println("Created user with id:", newUserId)

	return
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
