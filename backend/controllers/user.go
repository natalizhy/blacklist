package controllers

import (
	"bytes"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/natalizhy/blacklist/backend/models"
	"github.com/natalizhy/blacklist/backend/repositories"
	"gopkg.in/go-playground/validator.v9"
	"html/template"
	"io"
	"net/http"
	"strconv"
)

type UserTemp struct {
	User   models.User
	IsEdit bool

	IsSaveOk bool

	PhoneError     string
	CountryError   string
	LastNameError  string
	FirstNameError string
	InfoError      string
	PhotoError     string
	SearchError    string
}

var validate *validator.Validate

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repositories.GetUsers()

	if err != nil {
		w.Write([]byte("Юзеры не найден"))
		return
	}

	RenderTempl(w, "templates/users-list.html", users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userIDstr := chi.URLParam(r, "userID")
	user := models.User{}
	userTemp := UserTemp{IsEdit: false, User: user}

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

func GetNewUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{ID: 0}
	userTemp := UserTemp{IsEdit: true, User: user}

	RenderTempl(w, "templates/profile.html", userTemp)
}

var userError = map[string]map[string]string{
	"FirstName": {
		"required": "Имя не правельно введенное",
		"alpha":    "Можна использовать только буквы",
	},
	"LastName": {
		"required": "Имя не правельно введенное",
		"alpha":    "Можна использовать только буквы",
	},
	"Country": {
		"required": "Имя не правельно введенное",
		"alpha":    "Можна использовать только буквы",
	},
	"Phone": {
		"required": "Имя не правельно введенное",
		"numeric":  "Можна использовать только цифры",
	},
	"Info": {
		"required": "Обязательное поле",
	},
	"Photo": {
		"file": "файл",
	},
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	userTemp := UserTemp{IsEdit: true, User: user}

	// добавить валидации какието сюда

	user.FirstName = r.FormValue("first-name")
	user.LastName = r.FormValue("last-name")
	user.Country = r.FormValue("country")
	user.Phone = r.FormValue("phone")
	user.Info = r.FormValue("info")

	validate = validator.New()

	err := validate.Struct(user)

	if err != nil {
		//validationErrors := err.(validator.ValidationErrors)
		//fmt.Println(validationErrors)

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println()
			fmt.Println(userError[err.Field()][err.Tag()])
		}
		RenderTempl(w, "templates/profile.html", userTemp)

		return
	}

	userTemp.User = user

	userID, err := repositories.AddUser(user)

	if err != nil {
		w.Write([]byte("Юзер не добавлен"))
		return
	}

	http.Redirect(w, r, "/customers/"+strconv.FormatInt(userID, 10), http.StatusTemporaryRedirect)
}

func GetUpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDstr := chi.URLParam(r, "userID")
	user := models.User{}
	userTemp := UserTemp{IsEdit: true, User: user}

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

func Search(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("search"))
	userSearch := r.FormValue("search")

	user, err := repositories.Search(userSearch)

	if err != nil {
		w.Write([]byte("Юзер не найден"))
		return
	}

	tmplData := struct {
		UserSearch string
		User  []models.User
	}{
		UserSearch: userSearch,
		User:  user,
	}

	RenderTempl(w, "templates/search.html", tmplData)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDstr := chi.URLParam(r, "userID")
	user := models.User{}
	userTemp := UserTemp{IsEdit: true, User: user, IsSaveOk: false}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)

	if err != nil {
		w.Write([]byte("Юзер не найден"))
		return
	}
	user.ID = userID
	user.FirstName = r.FormValue("first-name")
	user.LastName = r.FormValue("last-name")
	user.Country = r.FormValue("country")
	user.Phone = r.FormValue("phone")
	user.Info = r.FormValue("info")

	userTemp.User = user

	err = repositories.UpdateUser(user, userID)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	userTemp.User = user
	userTemp.IsSaveOk = true

	RenderTempl(w, "templates/profile.html", userTemp)
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
