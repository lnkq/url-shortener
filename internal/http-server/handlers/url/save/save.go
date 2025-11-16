package save

import (
	"errors"
	"log/slog"
	"net/http"

	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	resp "url-shortener/internal/pkg/api/response"
	"url-shortener/internal/storage"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL       string `json:"url" validate:"required,url"`
	ShortCode string `json:"short_code" validate:"omitempty"`
}

type Response struct {
	resp.Response
	ShortCode string `json:"short_code,omitempty"`
}

// TODO: move to config
const shortCodeLength = 6

type URLSaver interface {
	SaveURL(original_url string, short_code string) error
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		// TODO FIXME
		// log = log.With(
		// 	slog.String("op", op),
		// 	slog.String("request_id", middleware.GetReqID(r.Context())),
		// )

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("invalid request body"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("request validation failed", sl.Err(err))
			render.JSON(w, r, resp.Error("invalid request parameters"))
			return
		}

		// TODO FIXME: handle existing URLs when auto-generating short codes

		short_code := req.ShortCode
		if short_code == "" {
			short_code = random.NewRandomString(shortCodeLength)
		}

		err = urlSaver.SaveURL(req.URL, short_code)
		if errors.Is(err, storage.ErrShortCodeExists) {
			log.Info("url already exists with this short code", slog.String("short_code", short_code))
			render.JSON(w, r, resp.Error("short code already exists"))
			return
		}
		if err != nil {
			log.Error("failed to save url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal server error"))
			return
		}

		log.Info("url saved successfully", slog.String("short_code", short_code))

		render.JSON(w, r, &Response{
			Response:  resp.Ok(),
			ShortCode: short_code,
		})
	}
}
