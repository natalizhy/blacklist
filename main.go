package main

import (
	"bytes"
	"fmt"
	//"github.com/CloudyKit/jet"
	"io"
	"regexp"
	"html/template"
	//"fmt"
	"github.com/go-chi/chi"
	//"io"
	"net/http"

	//"bytes"
	//"net/http"
	//"io"
	//"github.com/CloudyKit/jet"
	//"time"
	//"github.com/iizotop/esc/tradinghero.org/profile/jwt"

)

type (
	SignupUser struct {
		Phone     string
		Country   string
		LastName  string
		FirstName string
		Info      string

		PhoneError     string
		CountryError   string
		LastNameError  string
		FirstNameError string
		InfoError      string
	}
	SearchUser struct {
		Search string

		SearchError string
	}
)

func main() {

	mux := chi.NewRouter()

	mux.Post("/customers", Signup)
	mux.Get("/customers", Signup)
	mux.Get("/search", Search)


	fileHandle := http.FileServer(http.Dir(".")).ServeHTTP

	mux.Get("/css/*", fileHandle)
	mux.Get("/img/*", fileHandle)
	err := http.ListenAndServe(":3004", mux)
	fmt.Println(err)

}

func Signup(w http.ResponseWriter, r *http.Request) {
	signupUser := SignupUser{}

	if r.Method == "GET" {

		tmpl, err := template.ParseFiles("index.html")
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

		name, _ := regexp.MatchString("^[A-Za-zА-Яа-я]{4,}$", signupUser.FirstName)
		phone, _ := regexp.MatchString("^[0-9+ ()]{10,}$", signupUser.Phone)
		lname, _ := regexp.MatchString("^[A-Za-zА-Яа-я]{4,}$", signupUser.LastName)
		info, _ := regexp.MatchString("^[A-Za-zА-Яа-я]{4,}$", signupUser.Info)

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
		if signupUser.Country == "" {
			signupUser.CountryError = "Empty country"
		}
		if signupUser.Info == "" {
			signupUser.InfoError = "Empty country"
		} else if !info {
			signupUser.InfoError = "Use only Russian or English."
		} else if info {
			signupUser.InfoError = ""
		}

		if !name || !lname || !phone || !info {

			tmpl, err := template.ParseFiles("index.html")
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

			tmpl, err := template.ParseFiles("success.html")
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

		tmpl, err := template.ParseFiles("search.html")
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