package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"url-shortener2/internal/config"
	"url-shortener2/internal/http-server/handlers/url/redirect"
	"url-shortener2/internal/http-server/handlers/url/remove"
	"url-shortener2/internal/http-server/handlers/url/save"
	"url-shortener2/internal/http-server/middleware/httpLogger"
	"url-shortener2/internal/logging"
	"url-shortener2/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	//TODO: нужен конфиг: cleanenv
	cnfg := config.MustLoad()

	fmt.Printf("%#v\n", cnfg)
	//TODO: нужен логгер: slog
	logger := logging.SetupLogger(cnfg.Env)

	logger.Info("starting url-shortener", slog.String("env", cnfg.Env))
	logger.Debug("debug messages are enabled", slog.String("env", cnfg.Env))
	//TODO: инициируем storage: sqlite
	storage, err := sqlite.New(cnfg.StoragePath)
	if err != nil {
		logger.Info("failed to create storage: ", logging.Err(err))
		os.Exit(1)
	}

	//TODO: инициируем роутер: chi, render (chi render)
	router := chi.NewRouter()
	//middleware
	router.Use(middleware.RequestID)   //Присваиваем каждому запросу ID
	router.Use(middleware.RealIP)      //Чтобы посмотреть IP входящего запроса
	router.Use(httpLogger.New(logger)) //Используем написанный нами логгер для логирования обработки запросов
	router.Use(middleware.Recoverer)   //Восстанавливаемся после паники, которая может возникнуть внутри хэндлера
	router.Use(middleware.URLFormat)   //Позволяет парсить параметры URL

	// Создаём роутер для авторизации - Мы привносим в наш роутер новую логику для авторизации
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cnfg.ServerHTTP.User: cnfg.ServerHTTP.Password,
		}))

		// Роутеры перенесены сюда, так как требуют авторизации
		r.Post("/", save.New(logger, storage))
		r.Delete("/remove/{alias}", remove.New(logger, storage))
	})
	//Подключаем хэндлекры к роутеру
	router.Get("/{alias}", redirect.New(logger, storage))
	//TODO: запускаем сервер
	logger.Info("starting server", slog.String("address", cnfg.Address))
	server := &http.Server{
		Addr:         cnfg.Address,
		Handler:      router,
		ReadTimeout:  cnfg.ServerHTTP.Timeout,
		WriteTimeout: cnfg.ServerHTTP.Timeout,
		IdleTimeout:  cnfg.ServerHTTP.IdleTimeout,
	}
	defer server.Close()
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("failed to start server", logging.Err(err))
	}
	logger.Error("server stopped")
}
