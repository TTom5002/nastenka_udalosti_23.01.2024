package main

import (
	"nastenka_udalosti/internal/config"
	"nastenka_udalosti/internal/handlers"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)

	mux.Route("/user", func(mux chi.Router) {
		mux.Get("/login", handlers.Repo.Login)
		mux.Post("/login", handlers.Repo.PostLogin)

		mux.Get("/signup", handlers.Repo.Signup)
		mux.Post("/signup", handlers.Repo.PostSignup)

		mux.Get("/logout", handlers.Repo.Logout)
	})

	mux.Route("/dashboard", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Use(Verified)

		mux.Route("/cu", func(mux chi.Router) {
			mux.Route("/posts", func(mux chi.Router) {

				mux.Get("/make-event", handlers.Repo.MakeEvent)
				mux.Post("/make-event", handlers.Repo.PostMakeEvent)
				mux.Get("/my-events", handlers.Repo.MyEvents)

				mux.Get("/show-event/{id}", handlers.Repo.EditEvent)
				mux.Post("/show-event/{id}", handlers.Repo.PostUpdateEvent)
				mux.Get("/delete-event/{id}", handlers.Repo.DeleteEvent)
			})

			mux.Route("/profile", func(mux chi.Router) {
				mux.Get("/", handlers.Repo.EditProfile)
				mux.Post("/", handlers.Repo.PostEditProfile)
				mux.Get("/delete/{id}", handlers.Repo.DeleteProfile)
			})
		})

		mux.Route("/admin", func(mux chi.Router) {
			mux.Use(Admin)
			mux.Get("/", handlers.Repo.Home)
			mux.Get("/unverified-users", handlers.Repo.ShowAllUnverifiedUsers)
			mux.Post("/unverified-users", handlers.Repo.PostVerUsers)
			mux.Get("/all-users", handlers.Repo.AdminAllUsers)
			mux.Route("/posts", func(mux chi.Router) {
				mux.Get("/all-events", handlers.Repo.AdminShowAllEvents)

				// mux.Get("/", handlers.Repo.ViewProfile)
				// mux.Post("/edit", handlers.Repo.EditProfile)
				// Další akce s profilem
			})

		})

	})

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}