package main

import (
	"fmt"
	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
	"github.com/natalizhy/blacklist/backend/controllers"
	"github.com/natalizhy/blacklist/backend/repositories"
	"net/http"
)

func main() {

	repositories.InitDB()
	mux := chi.NewRouter()

	mux.Mount("/", adminRouter()) // проверка доступа

	mux.Get("/profiles/{userID}", controllers.GetUser) // просмотр юзера

	mux.Get("/", controllers.SearchForm)    // главная страница
	mux.Post("/", controllers.SearchResult) // поиск юзера

	fileHandle := http.FileServer(http.Dir(".")).ServeHTTP

	mux.Get("/assets/*", fileHandle)

	fmt.Println("Server was started at :3004")

	err := http.ListenAndServe(":3004", mux)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func adminRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Use(New("MyRealm", map[string][]string{
		"bob": {"password1", "password2"},
	}))

	mux.Post("/profiles/{userID}", controllers.GetUser)

	mux.Get("/profiles/{userID}/{mode}", controllers.GetUser) // редактирование
	mux.Post("/profiles/{userID}/{mode}", controllers.AddUser)

	mux.Get("/profiles/{userID}/DeleteUser", controllers.DeleteUser) // удаление юзера
	mux.Get("/profiles/{photoID}/DeletePhoto", controllers.DeleteUserPhoto) // удаление юзера

	mux.Get("/addNewUser", controllers.GetNewUser) //
	mux.Post("/addNewUser", controllers.AddUser)   // добавление нового профиля

	return mux
}

func New(realm string, credentials map[string][]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok {
				unauthorized(w, realm)
				return
			}

			validPasswords, userFound := credentials[username]
			if !userFound {
				unauthorized(w, realm)
				return
			}

			for _, validPassword := range validPasswords {
				if password == validPassword {
					next.ServeHTTP(w, r)
					return
				}
			}

			unauthorized(w, realm)
		})
	}
}

func unauthorized(w http.ResponseWriter, realm string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(http.StatusUnauthorized)
	http.Error(w, "Доступ ограничен", 401)
}
