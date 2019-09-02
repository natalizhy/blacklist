package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/natalizhy/blacklist/backend/models"
	"github.com/natalizhy/blacklist/backend/repositories"
	"github.com/natalizhy/blacklist/backend/utils"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type UserTemp struct {
	User     models.User
	Cities   map[int64]string
	Error    map[string]map[string]string
	IsEdit   bool
	IsSaveOk bool

	PhotoError string
}

type SearchUser struct {
	UserSearch string
	User       []models.User
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

func GetUser(w http.ResponseWriter, r *http.Request) {
	userIDstr := chi.URLParam(r, "userID")
	user := models.User{}
	userTemp := UserTemp{IsEdit: false, User: user, Cities: cities}

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
	RenderTempl(w, "templates/profile.html", userTemp)
}

func GetNewUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{ID: 0}
	userTemp := UserTemp{IsEdit: true, User: user, Cities: cities}

	RenderTempl(w, "templates/profile.html", userTemp)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	userIDstr := chi.URLParam(r, "userID")
	user := models.User{}
	userTemp := UserTemp{IsEdit: true, Cities: cities, PhotoError: ""}
	var err error
	var userID int64

	user.FirstName = r.FormValue("first-name")
	user.LastName = r.FormValue("last-name")
	user.CityID, _ = strconv.ParseInt(r.FormValue("city-id"), 10, 64)
	user.Phone = r.FormValue("phone")
	user.Info = r.FormValue("info")
	photo := r.FormValue("h-photo")

	file, _, photoErr := r.FormFile("photo")

	if photoErr != nil && photo == "" {
		userTemp.PhotoError = "Не выбрана фотография для юзера"
	}

	if photoErr == nil {
		defer file.Close()
	}

	userTemp.Error, err = utils.ValidateUser(user)

	userTemp.User = user

	if err == nil || userTemp.PhotoError == "" {

		if photoErr == nil {

			ct, err := getContentType(file)

			if err != nil {
				w.Write([]byte("err"))
				return
			}

			if _, ok := allowedMimeType[ct]; ok {
				user.Photo = allowedMimeType[ct]
			}
		}

		if photoErr != nil {
			user.Photo = photo
		}

		if userIDstr != "" {

			userID, err = strconv.ParseInt(userIDstr, 10, 64)
			if err != nil {
				w.Write([]byte("Неправельный ID"))
				return
			}

			err = repositories.UpdateUser(user, userID)

			if photo != "" {
				repositories.UpdateUserPhoto(user, userID)
			}

		} else {
			userID, err = repositories.AddUser(user)
			if err != nil {
				w.Write([]byte("Юзер не добавлен"))
				return
			}
		}

		if photoErr == nil {
			_ = SavePhoto(userID, file, user.Photo)
		}

		http.Redirect(w, r, "/profiles/"+strconv.FormatInt(userID, 10), http.StatusSeeOther)
		return
	}

	RenderTempl(w, "templates/profile.html", userTemp)
}

func getContentType(file multipart.File) (contentType string, err error) {
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println(23, err)
	}

	contentType = http.DetectContentType(buffer)

	_, err = file.Seek(0, 0)

	return
}

func SavePhoto(userID int64, file multipart.File, contentType string) (err error) {

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	err = ioutil.WriteFile("./assets/users-photo/"+strconv.FormatInt(userID, 10)+contentType, data, 0777)

	if err != nil {
		return
	}
	return
}

func GetUpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDstr := chi.URLParam(r, "userID")
	user := models.User{}
	userTemp := UserTemp{IsEdit: true, User: user, Cities: cities}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)

	if err != nil {
		w.Write([]byte("Юзер не найден"))
		return
	}

	user, err = repositories.GetUserById(userID)

	if err != nil {
		w.Write([]byte("Не могу выбрать юзера из базы"))
		return
	}

	userTemp.User = user

	RenderTempl(w, "templates/profile.html", userTemp)
}

func SearchGet(w http.ResponseWriter, r *http.Request) {
	userSearch := r.FormValue("search")
	user, err := repositories.Search(userSearch)
	tmplData := SearchUser{UserSearch: userSearch, User: user}

	if err != nil {
		w.Write([]byte("Юзеры не найден"))
		return
	}

	RenderTempl(w, "templates/users-list.html", tmplData)
}

func Search(w http.ResponseWriter, r *http.Request) {
	userSearch := r.FormValue("search")
	user, err := repositories.Search(userSearch)
	response := r.FormValue("g-recaptcha-response")
	tmplData := SearchUser{UserSearch: userSearch, User: user}

	fmt.Println("g-recaptcha-response : ", response)

	if response == "" {
		http.Redirect(w, r, "/", 301)
		return
	}

	remoteip := "176.38.148.28"

	fmt.Println("remote ip : ", remoteip)

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

	json.Unmarshal(body, &APIResp)
	fmt.Println(APIResp)

	if err != nil {
		w.Write([]byte("Юзер не найден"))
		return
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
