package repositories

import (
	"fmt"
	"github.com/natalizhy/blacklist/backend/models"
)

func GetUserById(userID int64) (user models.User, err error) {
	result := DB.QueryRowx("SELECT `id`, `first_name`, `last_name`, `city_id`, `phone`, `info` "+
		"FROM `profiles` WHERE `id`=? AND `status`=?", userID, 1)

	err = result.StructScan(&user)
	return
}

func GetPhotoById(userID int64) (photos []models.Photo, err error) {
	result, err := DB.Queryx("SELECT `id`, `userID`, `linkPhoto` FROM `photos` "+
		"WHERE `userID`=? AND `status`=?", userID, 1)

	if err != nil {
		return
	}

	for result.Next() {

		photo := models.Photo{}
		_ = result.StructScan(&photo)

		photos = append(photos, photo)
	}

	return
}

func UpdateUser(user models.User, userID int64) (err error) {
	_, err = DB.Exec("UPDATE `profiles` SET `first_name`=?, `last_name`=?, `city_id`=?, `phone`=?, `info`=? "+
		"WHERE id=?", user.FirstName, user.LastName, user.CityID, user.Phone, user.Info, userID)
	return
}

func AddUser(user models.User) (newUserId int64, err error) {
	newUserId = 0

	res, err := DB.Exec("INSERT INTO `profiles` (`first_name`, `last_name`, `city_id`, `phone`, `info`, `photoID`, `status`) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?)", user.FirstName, user.LastName, user.CityID, user.Phone, user.Info, user.PhotoID, 1)
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

func UpdateUserPhoto(user models.User, userID int64) (err error) {
	_, err = DB.Exec("UPDATE `profiles` SET `photoID`=? "+
		"WHERE id=?", user.PhotoID, userID)
	return
}

func AddUserPhoto(photo models.Photo, newUserId int64) (newUserPhotoId int64, err error) {
	newUserPhotoId = 0

	resp, err := DB.Exec("INSERT INTO `photos` (`userID`, `linkPhoto`, `status`) "+
		"VALUES (?, ?, ?)", newUserId, photo.LinkPhoto, 1)
	if err != nil {
		return
	}

	newUserPhotoId, err = resp.LastInsertId()
	if err != nil {
		return
	}

	fmt.Println("Created user with id:", newUserId, "and", photo.LinkPhoto, "linkPhoto")
	return
}

func Search(user string) (users []models.User, err error) {
	rows, err := DB.Queryx("SELECT p.`id`, p.`first_name`, p.`last_name`, p.`phone`, p.`photoID`, f.`linkPhoto` "+
		"FROM profiles AS p, photos AS f "+
		"WHERE p.status=1 AND p.photoID=f.id "+
		"AND ("+
		"p.first_name LIKE concat('%', ?, '%') OR "+
		"p.last_name LIKE concat('%', ?, '%') OR "+
		"p.city_id LIKE concat('%', ?, '%')"+
		")", user, user, user)

	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		user := models.User{}
		err = rows.StructScan(&user)
		if err != nil {
			fmt.Println(err)
			return
		}
		users = append(users, user)
	}

	return
}

func DeleteUser(user models.User, userID int64) (err error) {
	result := DB.QueryRowx("UPDATE `profiles` SET status=0 WHERE id=?", userID)

	err = result.StructScan(&user)

	fmt.Println(&user)
	fmt.Println("Delete user")

	return
}

func DeleteUserPhoto(photo models.Photo, photoID int64) (err error) {
	res := DB.QueryRowx("UPDATE `photos` SET status=0 WHERE id=?", photoID)

	err = res.StructScan(&photo)

	fmt.Println(&photo)
	fmt.Println("Delete photo user")

	return
}

func GetUserID(photoID int64) (UserID int64, err error) {
	result := DB.QueryRow("SELECT `userID` FROM `photos` "+
		"WHERE `ID`=?", photoID)

	err = result.Scan(&UserID)

	return
}
