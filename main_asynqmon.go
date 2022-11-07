package main

import (
	"log"
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
)

func main() {
	h := asynqmon.New(asynqmon.Options{
		RootPath: "/monitoring", // RootPath specifies the root for asynqmon app
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     "127.0.0.1:6379",
			DB:       2,
			Password: "Dc!123"},
	})

	// Note: We need the tailing slash when using net/http.ServeMux.
	http.Handle(h.RootPath()+"/", h)

	// Go to http://localhost:8080/monitoring to see asynqmon homepage.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
