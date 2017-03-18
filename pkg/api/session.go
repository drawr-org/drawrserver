package api

import (
	"context"
	"net/http"

	"github.com/drawr-team/drawrserver/pkg/bolt"
	"github.com/drawr-team/drawrserver/pkg/canvas"
	"github.com/drawr-team/drawrserver/pkg/session"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
)

const sessionIDParam string = "sessionID"

// sessionRouter sets up the session subroute
func sessionRouter() http.Handler {
	session.Init(dbClient)

	r := chi.NewRouter()
	notAllowed := r.MethodNotAllowedHandler()
	r.Get("/", sessionList)
	r.Put("/", sessionNewPUT)
	r.Post("/", notAllowed)
	r.Get("/new", sessionNewGET)
	r.Delete("/", notAllowed)

	r.Route("/:"+sessionIDParam, func(r chi.Router) {
		r.Use(sessionCtx)
		r.Get("/", sessionGet)
		r.Put("/", notAllowed)
		r.Post("/", sessionUpdate)
		r.Delete("/", sessionDelete)

		r.Get("/ws", sessionJoin)
		r.Get("/leave", sessionLeave)

	})

	return r
}

// WithSessionContext puts session into the context
func WithSessionContext(r *http.Request, s session.Session) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), sessionCtxKey, s))
}

// FromSessionContext gets session from context
func FromSessionContext(ctx context.Context) session.Session {
	return ctx.Value(sessionCtxKey).(session.Session)
}

func sessionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, sessionIDParam)
		s, err := session.Get(id)
		if err != nil {
			switch err {
			case bolt.ErrNotFound:
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		next.ServeHTTP(w, WithSessionContext(r, s))
	})
}

func sessionList(w http.ResponseWriter, r *http.Request) {
	sl, err := session.List()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
	}
	render.Status(r, http.StatusOK)
	wrapJSON(w, r, "sessions", sl)
}

func sessionNewPUT(w http.ResponseWriter, r *http.Request) {
	var data struct {
		session.Session
		OmitID interface{} `json:"id,omitempty"`
	}
	if err := render.Bind(r.Body, &data); err != nil {
		render.Status(r, http.StatusNotAcceptable)
		render.JSON(w, r, err.Error())
	}
	s, err := session.New(&data.Session)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, s)
}

func sessionNewGET(w http.ResponseWriter, r *http.Request) {
	s, err := session.New(nil)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, s)
}

func sessionGet(w http.ResponseWriter, r *http.Request) {
	s := FromSessionContext(r.Context())
	render.Status(r, http.StatusOK)
	render.JSON(w, r, s)
}

func sessionUpdate(w http.ResponseWriter, r *http.Request) {
	s := FromSessionContext(r.Context())
	var data struct {
		session.Session
		OmitID interface{} `json:"id,omitempty"`
	}
	if err := render.Bind(r.Body, &data); err != nil {
		render.Status(r, http.StatusNotAcceptable)
		render.JSON(w, r, err.Error())
	}
	if err := session.Update(s.ID, data.Session); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, data.Session)
}

func sessionDelete(w http.ResponseWriter, r *http.Request) {
	s := FromSessionContext(r.Context())
	if err := session.Delete(s); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, http.StatusText(http.StatusOK))
}

func sessionJoin(w http.ResponseWriter, r *http.Request) {
	s := FromSessionContext(r.Context())
	apilog.Println(s)
	if err := canvas.Connect(w, r, s); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, http.StatusText(http.StatusOK))
}

func sessionLeave(w http.ResponseWriter, r *http.Request) {
	s := FromSessionContext(r.Context())
	apilog.Println(s)
	// if err := service.Leave(session); err != nil {
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// }
	render.Status(r, http.StatusNotImplemented)
	render.PlainText(w, r, http.StatusText(http.StatusNotImplemented))
}
