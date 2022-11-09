package main

import (
	"log"
	"net/http"
)

type Myhandler struct {
}

func (ths *Myhandler) ServeHTTP(Writer http.ResponseWriter, Request *http.Request) {
	Writer.Write([]byte("myhandler"))
}

func main() {
	myMux := http.NewServeMux()
	myMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("NewServeMux"))
	})

	myMux.HandleFunc("/login", func(writer http.ResponseWriter, request *http.Request) {
		userName := request.URL.Query().Get("UserName")
		if userName != "" {
			c := http.Cookie{Name: "userName", Value: userName, Path: "/"}
			http.SetCookie(writer, &c)

		}
		writer.Write([]byte("login.html"))
	})

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello"))
	})
	http.HandleFunc("/abc", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("abc"))
	})
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
