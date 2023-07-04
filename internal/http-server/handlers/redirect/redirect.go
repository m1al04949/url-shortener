package redirect

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

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "hadlers.url.redirect.New"

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

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url nor found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", logslog.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got url", slog.String("url", resURL))

		// redirect to found URL
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
