package controllers

import (
	"bytes"
	"github.com/natalizhy/blacklist/backend/models"
	"github.com/natalizhy/blacklist/backend/repositories"
	"html/template"
	"io"
	"net/http"
	"regexp"
)

type (
	SignupUser struct {
		Phone     string
		Country   string
		LastName  string
		FirstName string
		Info      string
		Photo     string

		PhoneError     string
		CountryError   string
		LastNameError  string
		FirstNameError string
		InfoError      string
		PhotoError     string
	}
	SearchUser struct {
		Search string

		SearchError string
	}
)

func Signup(w http.ResponseWriter, r *http.Request) {
	signupUser := SignupUser{}

	if r.Method == "GET" {

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		body := &bytes.Buffer{}

		err = tmpl.Execute(body, signupUser)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		w.Write(body.Bytes())
	}

	if r.Method == "POST" {

		signupUser.Phone = r.FormValue("phone")
		signupUser.Country = r.FormValue("country")
		signupUser.LastName = r.FormValue("lname")
		signupUser.FirstName = r.FormValue("name")
		signupUser.Info = r.FormValue("info")
		signupUser.Photo = r.FormValue("photo")

		name, _ := regexp.MatchString("^[A-Za-zА-Яа-я]{4,}$", signupUser.FirstName)
		phone, _ := regexp.MatchString("^[0-9+ ()]{10,}$", signupUser.Phone)
		lname, _ := regexp.MatchString("^[A-Za-zА-Яа-я]{4,}$", signupUser.LastName)
		info, _ := regexp.MatchString("[^`]{4,}$", signupUser.Info)

		if signupUser.FirstName == "" {
			signupUser.FirstNameError = "Empty first name"
		} else if !name {
			signupUser.FirstNameError = "Use only Russian or English."
		} else if name {
			signupUser.FirstNameError = ""
		}
		if signupUser.LastName == "" {
			signupUser.LastNameError = "Empty last name"
		} else if !lname {
			signupUser.LastNameError = "Use only Russian or English."
		} else if lname {
			signupUser.LastNameError = ""
		}
		if signupUser.Phone == "" {
			signupUser.PhoneError = "Empty phone"
		} else if !phone {
			signupUser.PhoneError = "Use only numbers"
		} else if phone {
			signupUser.PhoneError = ""
		}
		if signupUser.Info == "" {
			signupUser.InfoError = "Empty information"
		} else if !info {
			signupUser.InfoError = "Use only Russian or English."
		} else if info {
			signupUser.InfoError = ""
		}
		if signupUser.Photo == "" {
			signupUser.PhotoError = "Empty Photo"
		}

		if !name || !lname || !phone || !info {

			tmpl, err := template.ParseFiles("templates/index.html")
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}

			body := &bytes.Buffer{}

			err = tmpl.Execute(body, signupUser)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}

			w.Write(body.Bytes())

		} else {

			tmpl, err := template.ParseFiles("templates/success.html")
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}

			body := &bytes.Buffer{}

			err = tmpl.Execute(body, signupUser)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}

			w.Write(body.Bytes())
		}
	}

	//w.Write([]byte("hello"))

}
func Search(w http.ResponseWriter, r *http.Request) {
	searchUser := SearchUser{}

	if r.Method == "GET" {

		tmpl, err := template.ParseFiles("templates/search.html")
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		body := &bytes.Buffer{}

		err = tmpl.Execute(body, searchUser)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		w.Write(body.Bytes())
	}

}

func AddNewUser(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		var (
			name= r.FormValue("name")
			lname= r.FormValue("lname")
			country= r.FormValue("country")
			phone= r.FormValue("phone")
			info= r.FormValue("info")
		)

		user := models.User{0, name, lname, country, phone, info}

		repositories.DBUser(user)
}
