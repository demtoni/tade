package api

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	StatusCode int    `json:"-"`
	ErrorMsg   string `json:"error"`
}

const (
	ErrorBadPassword          = "password should be at least 8 characters long and less than 72 characters long"
	ErrorPasswordsDontMatch   = "old and new passwords don't match."
	ErrorBadUsername          = "username should only contain numbers, letters or underscores and it's length should be less than 32 characters long."
	ErrorUsernameTaken        = "this username is already taken."
	ErrorInternal             = "internal server error, try again later."
	ErrorUnauthorized         = "you are not authorized to make this request."
	ErrorUserNotFound         = "user with this combination of username and password doesn't exist."
	ErrorBadInvite            = "this invite couldn't be used for registration."
	ErrorEmptyField           = "fields can't be empty."
	ErrorBadAmount            = "amount could not be negative and must be integer."
	ErrorServiceUnknown       = "unknown service type."
	ErrorServiceNotFound      = "service with that id doesn't exist."
	ErrorLocationNotSupported = "location doesn't support this type of service."
	ErrorLongServiceName      = "service name is too long (max 72 char)."
	ErrorNegativePeriod       = "service period must be > 0."
	ErrorLowBalance           = "your balance is too low."
	ErrorNoTransactions       = "no transactions found."
	ErrorNoServices           = "no services found."
	ErrorNoInvites            = "you can't generate any invites."
	ErrorNoUnusedInvites      = "you don't have any unused invites."
)

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func (s *Server) newErrorResponse(code int, msg string) render.Renderer {
	return &ErrorResponse{code, msg}
}

func (s *Server) SendError(w http.ResponseWriter, r *http.Request, err error, code int, msg string) {
	if err != nil {
		log.Println(err)
	}
	render.Render(w, r, s.newErrorResponse(code, msg))
}
