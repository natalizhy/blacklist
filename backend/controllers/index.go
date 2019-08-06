package controllers

import "net/http"

func IndexPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Index page"))
}