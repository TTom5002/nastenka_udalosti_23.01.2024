package repository

import "nastenka_udalosti/internal/models"

// DatabaseRepo uchovává postgres funkce
type DatabaseRepo interface {
	InsertEvent(event models.Event) error
	Authenticate(email, testPassword string) (models.User, error)
	ShowEvents() ([]models.Event, error)
	GetEventByID(event_id int) (models.Event, error)
	ShowUserEvents(id int) ([]models.Event, error)
	SignUpUser(user models.User) error
	DeleteEventByID(eventID int) error
	UpdateEventByID(event models.Event) error
	ShowUnverifiedUsers() ([]models.User, error)
	UpdateProfile(user models.User) error
	DeleteUserByID(userID int) error
	VerUserByID(userID int) error
	ShowAllUsers() ([]models.User, error)
	GetUserByID(ID int) (models.User, error)
}
