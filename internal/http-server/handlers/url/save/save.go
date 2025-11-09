package save

import (
	"errors"
	"net/http"
	"url-shortener2/internal/logging"
	"url-shortener2/internal/random"
	"url-shortener2/internal/response"
	"url-shortener2/internal/storage"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const aliasLength = 6

// URLSaver - интерфейс, который описывает объект, умеющий сохранять URL в базу данных
//go:generate mockery --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

// Мы создаём новый хэндлер, который умеет сохранять данные в базу и логировать
func New(logger *slog.Logger, db URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const origin = "http-server.handlers.url.save.New"

		logger = logger.With(
			slog.String("origin", origin),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req response.Request //Создаём объект запроса для парсинга

		err := render.DecodeJSON(r.Body, &req) //Парсим тело запроса с помощью декодера пакета chi/render
		if err != nil {
			logger.Error("failed to decode request body", logging.Err(err))

			render.JSON(w, r, response.Error("failed to decode request body")) //Пишем ответ пользователю с ошибкой в теле ответа с помощью пакета chi/render

			return
		}

		logger.Info("decoded request body", slog.Any("request", req))

		//Валидируем запрос  с помощью пакета github.com/go-playground/validator
		if err = validator.New().Struct(req); err != nil {
			logger.Error("failed to validate the request structure", logging.Err(err))

			render.JSON(w, r, response.ValidationError(err.(validator.ValidationErrors))) //возвращаем ответ, ошибки валидации в котором приведены в читаемый формат

			return
		}

		// Если алиас пустой, то генерируем его
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := db.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrUrlExists) {
				logger.Info("url already exists", slog.String("url", req.URL))
				render.JSON(w, r, response.Error("url alredy exists"))
				return
			}
			logger.Error("failed to save URL", logging.Err(err))

			render.JSON(w, r, response.Error("failed to save URL"))

			return
		}

		logger.Info("url added", slog.Int64("id", id))
		render.JSON(w, r, response.OK(alias))

	}
}
