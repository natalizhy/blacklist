package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"

	"github.com/natalizhy/blacklist/backend/controllers"
)

func main() {

	mux := chi.NewRouter()

	mux.Get("/", controllers.IndexPage)
	mux.Post("/customers", controllers.Signup)
	mux.Get("/customers", controllers.Signup)
	mux.Get("/search", controllers.Search)

	fileHandle := http.FileServer(http.Dir(".")).ServeHTTP

	mux.Get("/assets/*", fileHandle)

	fmt.Println("Server was started at :3004")

	err := http.ListenAndServe(":3004", mux)

	if err != nil {
		fmt.Println(err)
		return
	}

}
