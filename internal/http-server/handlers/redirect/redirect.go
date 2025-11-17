package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"url-shortener/internal/lib/logger/sl"
	resp "url-shortener/internal/pkg/api/response"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(short_code string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		shortCode := chi.URLParam(r, "short_code")
		if shortCode == "" {
			log.Info("short_code parameter is missing")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		orinigalURL, err := urlGetter.GetURL(shortCode)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("short_code does not match any URL", slog.String("short_code", shortCode))
			render.JSON(w, r, resp.Error("not found")) // ?
			return
		}
		if err != nil {
			log.Error("failed to get original URL", slog.String("short_code", shortCode), sl.Err(err))
			render.JSON(w, r, resp.Error("internal server error"))
			return
		}

		log.Info("redirecting to original URL",
			slog.String("short_code", shortCode),
			slog.String("original_url", orinigalURL),
		)

		http.Redirect(w, r, orinigalURL, http.StatusFound)
	}
}
