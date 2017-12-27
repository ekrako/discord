package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"github.com/ekrako/discord/timer"
)

func main() {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})

	// RESTy routes for "articles" resource
	r.Route("/timers", func(r chi.Router) {
		r.Get("/", ListTimers)
		r.Post("/", CreateTimer)
		r.Route("/{timerID}", func(r chi.Router) {
			r.Use(TimerCtx)
			r.Get("/", GetTimer)
			r.Put("/", UpdateTimer)
			r.Delete("/", DeleteTimer)
			r.Post("/start", StartTimer)
			r.Delete("/start", StopTimer)
			r.Post("/stop", StopTimer)
			r.Delete("/stop", StartTimer)
		})

	})
	log.Fatal(http.ListenAndServe(":3333", r))
}

// ListTimers List all the timers in server
func ListTimers(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, timer.GetAllTimers())
}

func TimerCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var thisTimer *timer.SingleTimer
		var err error

		if timerID := chi.URLParam(r, "timerID"); timerID != "" {
			thisTimer, err = timer.Get(timerID)
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "timer", thisTimer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CreateTimer Creates a new time
func CreateTimer(w http.ResponseWriter, r *http.Request) {
	data := &timer.Request{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	t, err := timer.Create(*data)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, t)
}

// GetTimer returns the specific Timer.
func GetTimer(w http.ResponseWriter, r *http.Request) {
	thisTimer := r.Context().Value("timer").(*timer.SingleTimer)
	render.JSON(w, r, thisTimer)
}

//StartTimer Make the timer running
func StartTimer(w http.ResponseWriter, r *http.Request) {
	thisTimer := r.Context().Value("timer").(*timer.SingleTimer)

	err := thisTimer.Start()
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.JSON(w, r, thisTimer)
}

//StopTimer Stop the timer running
func StopTimer(w http.ResponseWriter, r *http.Request) {
	thisTimer := r.Context().Value("timer").(*timer.SingleTimer)
	thisTimer.Stop()
	render.JSON(w, r, thisTimer)
}

// UpdateTimer the parameters of a specific Time specific Timer.
func UpdateTimer(w http.ResponseWriter, r *http.Request) {
	data := &timer.Request{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	thisTimer := r.Context().Value("timer").(*timer.SingleTimer)
	thisTimer.Update(*data)
	render.JSON(w, r, thisTimer)
}

// DeleteTimer an existing timer
func DeleteTimer(w http.ResponseWriter, r *http.Request) {
	thisTimer := r.Context().Value("timer").(*timer.SingleTimer)

	thisTimer.Delete()

	render.NoContent(w, r)
}

// Binder interface for managing request payloads.
type Binder interface {
	Bind(r *http.Request) error
}

// Bind decodes a request body and executes the Binder method of the
// payload structure.
func Bind(r *http.Request, v Binder) error {
	if err := render.Decode(r, v); err != nil {
		return err
	}
	return render.Bind(r, v)
}

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
