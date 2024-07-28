package main

import (
	"fmt"
	"log"
	"net/http"
	"skillfactory/GoNews/pkg/api"
	"skillfactory/GoNews/pkg/storage"

	//database "skillfactory/GoNews/pkg/storage/memdb"
	database "skillfactory/GoNews/pkg/storage/postgres"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server

	pwd := "passwod_to_database"
	host := "host_database"

	sqlconn := fmt.Sprintf("postgres://postgres:%s@%s/posts", pwd, host)
	// Создаём объекты баз данных.
	//
	/*	// БД в памяти.
		db, err := database.New("")
	*/

	db, err := database.New(sqlconn)
	if err != nil {
		log.Fatal(err)
	}
	/*
		// Документная БД MongoDB.
			db3, err := mongo.New("mongodb://server.domain:27017/")
			if err != nil {
				log.Fatal(err)
			}
			_, _ = db2, db3
	*/

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	http.ListenAndServe(":8081", srv.api.Router())
}
