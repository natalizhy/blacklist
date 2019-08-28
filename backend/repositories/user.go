package repositories

import (
	"fmt"
	"github.com/natalizhy/blacklist/backend/models"
)

func GetUsers() (users []models.User, err error) {
	rows, err := DB.Queryx("SELECT `id`, `first_name`, `last_name`, `city_id`, `phone`, `info`, `photo` " +
		"FROM `profiles` " +
		"WHERE `status`=1 " +
		"ORDER BY id DESC")

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
	result := DB.QueryRowx("SELECT `id`, `first_name`, `last_name`, `city_id`, `phone`, `info`, `photo` "+
		"FROM `profiles` WHERE `id`=?", userID)

	err = result.StructScan(&user)

	return
}

func UpdateUser(user models.User, userID int64) (err error) {
	_, err = DB.Exec("UPDATE `profiles` SET `first_name`=?, `last_name`=?, `city_id`=?, `phone`=?, `info`=? "+
		"WHERE id=?", user.FirstName, user.LastName, user.CityID, user.Phone, user.Info, userID)
	return
}

func UpdateUserPhoto(user models.User, userID int64) (err error) {
	_, err = DB.Exec("UPDATE `profiles` SET `photo`=? "+
		"WHERE id=?", user.Photo, userID)
	return
}

func AddUser(user models.User) (newUserId int64, err error) {
	newUserId = 0

	res, err := DB.Exec("INSERT INTO `profiles` (`first_name`, `last_name`, `city_id`, `phone`, `info`, `photo`, `status`) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?)", user.FirstName, user.LastName, user.CityID, user.Phone, user.Info, user.Photo, 1)
	if err != nil {
		return
	}
	newUserId, err = res.LastInsertId()
	if err != nil {
		return
	}

	fmt.Println("Created user with id:", newUserId, "and", user.Photo, "photo")

	return
}

func Search(user string) (users []models.User, err error) {
	rows, err := DB.Queryx("SELECT `id`, `first_name`, `last_name`, `city_id`, `photo` "+
		"FROM profiles "+
		"WHERE status=1 "+
		"AND ("+
		"first_name LIKE concat('%', ?, '%') OR "+
		"last_name LIKE concat('%', ?, '%') OR "+
		"city_id LIKE concat('%', ?, '%')"+
		")", user, user, user)

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
	result := DB.QueryRowx("UPDATE `profiles` SET status=0 WHERE id=?", userID)

	err = result.StructScan(&user)

	fmt.Println("Delete user")

	return
}
