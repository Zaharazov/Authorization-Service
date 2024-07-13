package app

import (
	"Authorization-Service/internal/presentation/routers"
	"Authorization-Service/internal/server/configs"
	"Authorization-Service/internal/server/redis"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Run() {
	log.Printf("Server started")

	log.Println(redis.Client)

	router := routers.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL("http://localhost" + configs.Port + "/swagger/doc.json")))

	log.Fatal(http.ListenAndServe(configs.Port, router))
}
