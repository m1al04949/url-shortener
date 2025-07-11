package delete

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/m1al04949/url-shortener/internal/lib/api/response"
	"github.com/m1al04949/url-shortener/internal/lib/logger/logslog"
	"github.com/m1al04949/url-shortener/internal/storage"
	"golang.org/x/exp/slog"
)

type Response struct {
	Response resp.Response
	Alias    string `json:"alias"`
}

type URLDeleter interface {
	DeleteURL(alias string) error
}

// DeleteShortURL godoc
// @Summary Удалять короткую ссылку
// @Description Удаляет короткий URL
// @Tags url
// @Accept  json
// @Produce  json
// @Param alias path string true "Алиас короткой ссылки"
// @Success 200 {object} Response "Successfully deleted"
// @Failure 400 {object} resp.Response "Invalid request"
// @Failure 404 {object} resp.Response "URL not found"
// @Failure 500 {object} resp.Response "Internal server error"
// @Router /url/{alias} [delete]
func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to delete url", logslog.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("url is deleted by", slog.String("alias", alias))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
