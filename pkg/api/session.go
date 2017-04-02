package api

import (
	"context"
	"net/http"

	"github.com/drawr-team/drawrserver/pkg/bolt"
	"github.com/drawr-team/drawrserver/pkg/canvas"
	"github.com/drawr-team/drawrserver/pkg/service"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
)

const sessionIDParam string = "sessionID"

type sessionResponse struct {
	service.Session
	OmitID interface{} `json:"id,omitempty"`
}

func (sr *sessionResponse) Bind(r *http.Request) error {
	return nil
}

// sessionRouter sets up the session subroute
func sessionRouter() http.Handler {
	service.Init(dbClient)

	r := chi.NewRouter()
	notAllowed := r.MethodNotAllowedHandler()
	r.Get("/", sessionList)
	r.Put("/", sessionNewPUT)
	r.Post("/", notAllowed)
	r.Get("/new", sessionNewGET)
	r.Delete("/", notAllowed)

	// reroute old test endpoint
	r.Get("/__test__", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/version", http.StatusMovedPermanently)
	})

	r.Route("/:"+sessionIDParam, func(r chi.Router) {
		r.Use(sessionCtx)
		r.Get("/", sessionGet)
		r.Put("/", notAllowed)
		r.Post("/", sessionUpdate)
		r.Delete("/", sessionDelete)

		r.Get("/ws", sessionJoin)
		r.Get("/leave", sessionLeave)

		r.Mount("/user", userRouter())
	})

	return r
}

func withSessionContext(r *http.Request, s service.Session) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), sessionCtxKey, s))
}

func fromSessionContext(ctx context.Context) service.Session {
	return ctx.Value(sessionCtxKey).(service.Session)
}

func sessionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, sessionIDParam)
		s, err := service.GetSession(id)
		if err != nil {
			switch err {
			case bolt.ErrNotFound:
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		next.ServeHTTP(w, withSessionContext(r, s))
	})
}

func sessionList(w http.ResponseWriter, r *http.Request) {
	sl, err := service.ListSessions()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
	}
	render.Status(r, http.StatusOK)
	wrapJSON(w, r, "sessions", sl)
}

func sessionNewPUT(w http.ResponseWriter, r *http.Request) {
	var data sessionResponse
	if err := render.Bind(r, &data); err != nil {
		render.Status(r, http.StatusNotAcceptable)
		render.JSON(w, r, err.Error())
	}
	s, err := service.NewSession(&data.Session)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, s)
}

func sessionNewGET(w http.ResponseWriter, r *http.Request) {
	s, err := service.NewSession(nil)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, s)
}

func sessionGet(w http.ResponseWriter, r *http.Request) {
	s := fromSessionContext(r.Context())
	render.Status(r, http.StatusOK)
	render.JSON(w, r, s)
}

func sessionUpdate(w http.ResponseWriter, r *http.Request) {
	s := fromSessionContext(r.Context())
	var data sessionResponse
	if err := render.Bind(r, &data); err != nil {
		render.Status(r, http.StatusNotAcceptable)
		render.JSON(w, r, err.Error())
	}
	if err := service.UpdateSession(s.ID, data.Session); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, data.Session)
}

func sessionDelete(w http.ResponseWriter, r *http.Request) {
	s := fromSessionContext(r.Context())
	if err := service.DeleteSession(s); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, http.StatusText(http.StatusOK))
}

func sessionJoin(w http.ResponseWriter, r *http.Request) {
	s := fromSessionContext(r.Context())
	apilog.Println(s)
	if err := canvas.Connect(w, r, s); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err.Error())
	}
}

func sessionLeave(w http.ResponseWriter, r *http.Request) {
	s := fromSessionContext(r.Context())
	apilog.Println(s)
	// if err := service.Leave(session); err != nil {
	//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// }
	render.Status(r, http.StatusNotImplemented)
	render.PlainText(w, r, http.StatusText(http.StatusNotImplemented))
}
