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

	mux.Get("/", controllers.IndexPage)
	mux.Post("/customers", controllers.Signup)
	mux.Get("/customers", controllers.Signup)
	mux.Get("/search", controllers.Search)
	mux.Post("/", controllers.AddNewUser)

	fileHandle := http.FileServer(http.Dir(".")).ServeHTTP

	mux.Get("/assets/*", fileHandle)

	fmt.Println("Server was started at :3004")

	err := http.ListenAndServe(":3004", mux)

	if err != nil {
		fmt.Println(err)
		return
	}

}
