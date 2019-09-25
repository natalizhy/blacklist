package controllers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/natalizhy/blacklist/backend/models"
	"github.com/natalizhy/blacklist/backend/repositories"
	"github.com/natalizhy/blacklist/backend/utils"
	"github.com/nfnt/resize"
	"html/template"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type UserTemp struct {
	User     models.User
	Cities   map[int64]string
	Error    map[string]map[string]string
	IsEdit   bool
	IsSaveOk bool

	ListPhotos []models.Photo
	Photo      models.Photo
	PhotoError string
}

type SearchUser struct {
	UserSearch   string
	User         []models.User
	ReCAPTCHAerr string
	Mismatch     string
}

var cities = map[int64]string{
	1: "Киев",
	2: "Харков",
	3: "Одесса",
}

var allowedMimeType = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

type JSONAPIResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []int     `json:"error-codes"`
}

var rand uint32
var randmu sync.Mutex

func GetUser(w http.ResponseWriter, r *http.Request) {
	userIDstr := chi.URLParam(r, "userID")

	user := models.User{}
	photo := models.Photo{}
	photos := []models.Photo{}

	userTemp := UserTemp{User: user, Cities: cities, Photo: photo}

	IsEdit := chi.URLParam(r, "mode") == "edit"
	if IsEdit == true {
		userTemp.IsEdit = true
	} else {
		userTemp.IsEdit = false
	}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		w.Write([]byte("Профиль не найден"))
		return
	}

	user, err = repositories.GetUserById(userID)
	if err != nil {
		w.Write([]byte("Не могу выбрать профиль из базы"))
		return
	}
	userTemp.User = user

	photos, err = repositories.GetPhotoById(userID)
	if err != nil {
		fmt.Println("err", err)
	}

	userTemp.ListPhotos = photos
	RenderTempl(w, "templates/profile.html", userTemp)

}

func GetNewUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{ID: 0}
	photo := models.Photo{ID: 0}
	userTemp := UserTemp{IsEdit: true, User: user, Cities: cities, Photo: photo}

	RenderTempl(w, "templates/profile.html", userTemp)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	userIDstr := chi.URLParam(r, "userID")
	user := models.User{}
	photo := models.Photo{}
	userTemp := UserTemp{IsEdit: true, Cities: cities, PhotoError: "", Photo: photo}
	var err error
	var userID int64
	var photoID int64

	user.FirstName = r.FormValue("first-name")
	user.LastName = r.FormValue("last-name")
	user.CityID, _ = strconv.ParseInt(r.FormValue("city-id"), 10, 64)
	user.Phone = r.FormValue("phone")
	user.Info = r.FormValue("info")

	files := r.MultipartForm.File["photo"]
	if files == nil {
		userTemp.PhotoError = "Не выбрана фотография"
	}

	if userIDstr != "" {
		userTemp.PhotoError = ""
	}

	userTemp.Error, err = utils.ValidateUser(user, photo)
	if err != nil {
		w.Write([]byte("err"))
		return
	}

	userTemp.User = user

	usph, err := AddUserPhoto(w, r, userTemp)
	if err != nil {
		w.Write([]byte("err"))
		return
	}

	if userTemp.PhotoError == "" {

		if userIDstr != "" {

			userID, err = strconv.ParseInt(userIDstr, 10, 64)
			if err != nil {
				w.Write([]byte("Неправельный ID"))
				return
			}

			err = repositories.UpdateUser(user, userID)
			if err != nil {
				fmt.Println(err)
				w.Write([]byte("Юзер не добавлен"))
				return
			}

		} else {
			userID, err = repositories.AddUser(user)

			if err != nil {
				fmt.Println(err)
				w.Write([]byte("Юзер не добавлен"))
				return
			}
		}

		photoID, err = RangeUserPhoto(w, photo, userID, usph)
		if err != nil {
			fmt.Println(err)
			w.Write([]byte("Главное фото для юзера не выбрано"))
			return
		}
		user.PhotoID = photoID

		photoIDup := repositories.UpdateUserPhoto(user, userID)
		if photoIDup != nil {
			fmt.Println(photoIDup)
		}

		fmt.Println(user.PhotoID, photoIDup, "done")

		http.Redirect(w, r, "/profiles/"+strconv.FormatInt(userID, 10), http.StatusSeeOther)
		return
	}

	RenderTempl(w, "templates/profile.html", userTemp)
}

func Search(w http.ResponseWriter, r *http.Request) {
	userSearch := r.FormValue("search")
	user, err := repositories.Search(userSearch)
	tmplData := SearchUser{UserSearch: userSearch, User: user, ReCAPTCHAerr: "Докажите, что вы не робот", Mismatch: ""}

	if len(user) == 0 {
		tmplData.Mismatch = "Ничего по вашему запросу не найдено, возможно в строке поиска вы допустили ошибку"
	}

	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Профиль не найден"))
		return
	}

	if userSearch != "" {
		response := r.FormValue("g-recaptcha-response")

		if response == "" {
			http.Redirect(w, r, "/", 301)
			return
		}

		remoteip := "176.38.148.28"

		secret := "6LcGMLYUAAAAAO7SPd_o6HjAqpHe_VH4CrX5kA3d"
		postURL := "https://www.google.com/recaptcha/api/siteverify"

		postStr := url.Values{"secret": {secret}, "response": {response}, "remoteip": {remoteip}}

		responsePost, err := http.PostForm(postURL, postStr)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer responsePost.Body.Close()

		body, err := ioutil.ReadAll(responsePost.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		var APIResp JSONAPIResponse

		err = json.Unmarshal(body, &APIResp)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/", 301)
			return
		}
	}

	RenderTempl(w, "templates/users-list.html", tmplData)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDstr := chi.URLParam(r, "userID")

	user := models.User{}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)

	err = repositories.DeleteUser(user, userID)

	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
}

func DeleteUserPhoto(w http.ResponseWriter, r *http.Request) {
	photoIDstr := chi.URLParam(r, "photoID")

	photo := models.Photo{}
	fmt.Println(photo.UserID)

	photoID, err := strconv.ParseInt(photoIDstr, 10, 64)
	fmt.Println(err, photoID, photo.UserID)

	userID, err := repositories.GetUserID(photoID)
	fmt.Println(err, userID, photo.UserID)

	err = repositories.DeleteUserPhoto(photo, photoID)

	if err != nil {
		http.Redirect(w, r, "/profiles/"+strconv.FormatInt(userID, 10), http.StatusTemporaryRedirect)
		return
	}
}

func RenderTempl(w http.ResponseWriter, tmplName string, data interface{}) {
	tmpl, err := template.ParseFiles(tmplName)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	body := &bytes.Buffer{}

	err = tmpl.Execute(body, data)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	w.Write(body.Bytes())
}

func Hash(path string) (string, error) {
	var returnMD5String string

	files, err := os.Open(path)
	if err != nil {
		return returnMD5String, err
	}

	defer files.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, files); err != nil {
		return returnMD5String, err
	}

	hashInBytes := hash.Sum(nil)[:16]

	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil
}

func AddUserPhoto(w http.ResponseWriter, r *http.Request, userTemp UserTemp) (list []string, err error) {
	files := r.MultipartForm.File["photo"]

	for i := range files {

		file, err := files[i].Open()
		fmt.Println(file, "file")
		if err != nil {
			fmt.Fprintln(w, err)
			userTemp.PhotoError = "Не выбрана фотография для юзера"
		}

		imgDecode, format, err := image.Decode(file)
		fmt.Println(file, format, err, "imgDecode")

		if err != nil {
			fmt.Println(err, "err")
			fmt.Fprintln(w, err)
			userTemp.PhotoError = "Не раскодировалось"
		}

		defer file.Close()

		//imgDecode = resize.Resize(220, 300, imgDecode, resize.Bicubic) // <-- изменение размера картинки
		imgDecode = resize.Thumbnail(220, 300, imgDecode, resize.Bicubic) // <-- изменение размера картинки

		tmpFileName := fmt.Sprintf("%x.tmp", time.Now().UnixNano())
		tmpFilePath := "./assets/users-photo/" + tmpFileName
		pathRes := "./assets/users-photo/"

		out, err := os.Create(tmpFilePath)
		if err != nil {
			fmt.Fprintln(w, err)
		}

		png.Encode(out, imgDecode)

		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		}

		ext := filepath.Ext(files[i].Filename)

		hash, err := Hash(tmpFilePath)
		if err != nil {
			w.Write([]byte("Ошибка с хешем"))
		}

		err = os.Rename(tmpFilePath, pathRes+hash+ext)
		if err != nil {
			fmt.Println(err)
		}

		list = append(list, pathRes+hash+ext)
	}
	return
}

func RangeUserPhoto(w http.ResponseWriter, photo models.Photo, userID int64, usph []string) (photoID int64, err error) {

	for _, val := range usph {

		res := val[1:]
		photo.LinkPhoto = res
		fmt.Println(photo.LinkPhoto, val, "val")

		photoID, err = repositories.AddUserPhoto(photo, userID)

		fmt.Println(userID, photoID, "ID2")
		if err != nil {
			fmt.Println(err)
			w.Write([]byte("Фотографии не добавлены"))
			return
		}
	}
	return
}
