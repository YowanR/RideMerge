package ridemerge

import (
	"log"

	"cloud.google.com/go/datastore"

	"golang.org/x/net/context"
)

var (
	PROJECT_ID string
	UDB        UserDatabase
	TDB        TripDatabase
)

func init() {
	var err error
	PROJECT_ID = "ridemerge"

	UDB, err = configureUserDatastore(PROJECT_ID)
	if err != nil {
		log.Fatal(err)
	}

	TDB, err = configureTripDatastore(PROJECT_ID)
	if err != nil {
		log.Fatal(err)
	}
}

func configureUserDatastore(projectID string) (UserDatabase, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return NewUserDatastore(client)
}

func configureTripDatastore(projectID string) (TripDatabase, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return NewTripDatastore(client)
}
