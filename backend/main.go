package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/natalizhy/blacklist/backend/controllers"
	"github.com/natalizhy/blacklist/backend/repositories"
)

var database *sql.DB

func main() {

	repositories.InitDB()
	mux := chi.NewRouter()

	mux.HandleFunc("/signup", use(myHandler, controllers.Signup))

	mux.Get("/", controllers.GetUsers)
	//mux.Post("/customers/{userID}", controllers.Redirect)
	mux.Get("/profiles/{userID}", controllers.GetUser) // просмотр юзера
	mux.Post("/profiles/{userID}", controllers.GetUser)
	mux.Get("/profiles/{userID}/edit", controllers.GetUpdateUser) // редактирование
	mux.Post("/profiles/{userID}/edit", controllers.AddUser)
	mux.Get("/profiles/{userID}/Delete", controllers.DeleteUser) // удаление юзера

	mux.Get("/signup", controllers.SignupForm) // Регистрация

	mux.Get("/addNewUser", controllers.GetNewUser) //
	mux.Post("/addNewUser", controllers.AddUser)   // добавление нового юзера

	mux.Get("/searchUser", controllers.Search)  //
	mux.Post("/searchUser", controllers.Search) // поиск юзера

	//mux.Post("/customers", controllers.AddNewUser)
	//mux.Post("/", controllers.AddNewUser)

	fileHandle := http.FileServer(http.Dir(".")).ServeHTTP

	mux.Get("/assets/*", fileHandle)

	fmt.Println("Server was started at :3004")

	err := http.ListenAndServe(":3004", mux)

	if err != nil {
		fmt.Println(err)
		return
	}

}

func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

func myHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Authenticated!"))
	return
}
