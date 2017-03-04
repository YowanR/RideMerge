package ridemerge

// User model
type User struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Trips     []Trip
}

// Session model
type Session struct {
	User
	LoggedIn      bool
	SearchResults []Trip
}

// RidePost model
type Trip struct {
	Event       string
	CreatedBy   string
	Origin      string
	Destination string
	DDate       string
	DTime       string
	HasCar      string
	Seats       int
}

type UserDatabase interface {
	ListUsers() ([]User, error)

	GetUser(email string) (User, error)

	PutUser(u User) error

	DeleteUser(email string) error

	UpdateUser(u User) error

	Close()
}

type TripDatabase interface {
	ListTrips() ([]Trip, error)

	ListTripsWithDestination(dest string) ([]Trip, error)

	ListTripsCreatedBy(createdBy string) ([]Trip, error)

	GetTrip(event string) (Trip, error)

	PutTrip(t Trip) error

	UpdateTrip(t Trip) error

	DeleteTrip(event string) error

	Close()
}
