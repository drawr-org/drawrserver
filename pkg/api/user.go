package api

import (
	"context"
	"net/http"

	"github.com/drawr-team/drawrserver/pkg/bolt"
	"github.com/drawr-team/drawrserver/pkg/service"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
)

const userIDParam string = "userID"

type userResponse struct {
	service.User
	OmitID interface{} `json:"id,omitempty"`
}

func (ur *userResponse) Bind(r *http.Request) error {
	return nil
}

func userRouter() http.Handler {
	service.Init(dbClient)

	r := chi.NewRouter()
	notAllowed := r.MethodNotAllowedHandler()

	// on user list
	r.Get("/", userList)
	r.Put("/", notAllowed)
	r.Post("/", notAllowed)
	r.Delete("/", notAllowed)
	// on single user
	r.Route("/:"+userIDParam, func(r chi.Router) {
		r.Use(userCtx)
		r.Get("/", userGet)
		r.Put("/", userNew)
		r.Post("/", userUpdate)
		r.Delete("/", userDelete)
	})

	return r
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotImplemented)
	render.PlainText(w, r, http.StatusText(http.StatusNotImplemented))
}

func withUserContext(r *http.Request, u service.User) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), userCtxKey, u))
}

func fromUserContext(ctx context.Context) service.User {
	return ctx.Value(userCtxKey).(service.User)
}

func userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, userIDParam)
		u, err := service.GetUser(id)
		if err != nil {
			switch err {
			case bolt.ErrNotFound:
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		next.ServeHTTP(w, withUserContext(r, u))
	})
}

func userList(w http.ResponseWriter, r *http.Request) {
	ul, err := service.ListUsers()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
	}
	render.Status(r, http.StatusOK)
	wrapJSON(w, r, "users", ul)
}

func userGet(w http.ResponseWriter, r *http.Request) {
	u := fromUserContext(r.Context())
	render.Status(r, http.StatusOK)
	render.JSON(w, r, u)
}

func userNew(w http.ResponseWriter, r *http.Request) {
	u := fromUserContext(r.Context())
	var data userResponse
	if err := render.Bind(r, &data); err != nil {
		render.Status(r, http.StatusNotAcceptable)
		render.JSON(w, r, err.Error())
	}
	u, err := service.NewUser(&data.User)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, u)
}

func userUpdate(w http.ResponseWriter, r *http.Request) {
	u := fromUserContext(r.Context())
	var data userResponse
	if err := render.Bind(r, &data); err != nil {
		render.Status(r, http.StatusNotAcceptable)
		render.JSON(w, r, err.Error())
	}
	if err := service.UpdateUser(u.ID, data.User); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, data.User)
}

func userDelete(w http.ResponseWriter, r *http.Request) {
	u := fromUserContext(r.Context())
	if err := service.DeleteUser(u); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, http.StatusText(http.StatusOK))
}
