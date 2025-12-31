package save

import (
	"errors"
	resp "finance1/internal/lib/api/response"
	"finance1/internal/lib/random"
	"finance1/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

// TODO move to conf
const aliasLength = 6

var validate = validator.New()

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", slog.Any("err", err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		if err := validate.Struct(req); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				render.JSON(w, r, resp.ValidationError(ve))
				return
			}
			log.Error("validator error", slog.Any("err", err))
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", req.URL))
				render.JSON(w, r, resp.Error("url already exists"))
				return
			}

			log.Error("failed to add url", slog.Any("err", err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}

		log.Info("url added", slog.Int64("id", id))
		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
