package repositories

import (
	//"database/sql"
	"fmt"
	"github.com/natalizhy/blacklist/backend/models"
	//"github.com/jmoiron/sqlx"
)

func GetUsers() (users []models.User, err error) {
	rows, err := DB.Queryx("SELECT `id`, `first_name`, `last_name`, `country`, `phone`, `info` FROM `customers`  WHERE `status`=1 ORDER BY id DESC")

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
	result := DB.QueryRowx("SELECT `id`, `first_name`, `last_name`, `country`, `phone`, `info` FROM `customers` WHERE `id`=?", userID)

	err = result.StructScan(&user)

	return
}

func UpdateUser(user models.User, userID int64) (err error) {
	_, err = DB.Exec("UPDATE `customers` SET `first_name`=?, `last_name`=?, `country`=?, `info`=?, `phone`=? WHERE id=?", user.FirstName, user.LastName, user.Country, user.Info, user.Phone, userID)
	return
}

func AddUser(user models.User) (newUserId int64, err error) {
	newUserId = 0

	res, err := DB.Exec("INSERT INTO `customers` (`first_name`, `last_name`, `country`, `phone`, `info`, `status`) VALUES (?, ?, ?, ?, ?, ?)", user.FirstName, user.LastName, user.Country, user.Phone, user.Info, 1)
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
func Search(user string) (users []models.User, err error) {
	rows, err := DB.Queryx("SELECT `last_name` FROM customers WHERE  `last_name` LIKE concat('%', ?, '%')", user)

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
func DeleteUser(user models.User, userID int64) (err error) {
	result := DB.QueryRowx("UPDATE `customers` SET status=? WHERE id=?", 0, userID)

	err = result.StructScan(&user)

	fmt.Println("Delete user")

	return
}
