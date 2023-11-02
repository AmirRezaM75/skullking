package routes

import (
	"github.com/AmirRezaM75/skull-king/contracts"
	"github.com/AmirRezaM75/skull-king/handlers"
	"github.com/AmirRezaM75/skull-king/middlewares"
	"github.com/amirrezam75/go-router"
)

type Route struct {
	Router      *router.Router
	UserService contracts.UserService
	UserHandler handlers.UserHandler
	GameHandler *handlers.GameHandler
}

func (r Route) Setup() {
	r.Router.Get("/verify-email/:id/:hash", r.UserHandler.VerifyEmail).
		Middleware(middlewares.ValidateSignature{})
	r.Router.Post("/register", r.UserHandler.Register)
	r.Router.Post("/login", r.UserHandler.Login)
	r.Router.Post("/email/verification-notification", r.UserHandler.EmailVerificationNotification).
		Middleware(middlewares.Authenticate{UserService: r.UserService})
	r.Router.Post("/forgot-password", r.UserHandler.ForgotPassword)
	r.Router.Post("/reset-password", r.UserHandler.ResetPassword)

	r.Router.Post("/games", r.GameHandler.Create).
		Middleware(middlewares.Authenticate{UserService: r.UserService})
	r.Router.Get("/games/join", r.GameHandler.Join)
	r.Router.Get("/games/cards", r.GameHandler.Cards)
}
