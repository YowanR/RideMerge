package ridemerge

import (
	"fmt"
	"log"
	"strings"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
)

type TripDatastore struct {
	client *datastore.Client
}

var _ TripDatabase = &TripDatastore{}

func NewTripDatastore(client *datastore.Client) (TripDatabase, error) {
	ctx := context.Background()
	t, err := client.NewTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("tripdatastore: could not connect: %v", err)
	}
	if err := t.Rollback(); err != nil {
		return nil, fmt.Errorf("tripdatastore: could not connect: %v", err)
	}
	return &TripDatastore{
		client: client,
	}, nil
}

// Close closes the database.
func (db *TripDatastore) Close() {
	// No op.
}

func (db *TripDatastore) datastoreKey(event string) *datastore.Key {
	return datastore.NameKey("Trip", event, nil)
}

func (db *TripDatastore) GetTrip(event string) (Trip, error) {
	ctx := context.Background()
	k := datastore.NameKey("Trip", event, nil)
	var trip Trip
	if err := db.client.Get(ctx, k, &trip); err != nil {
		return trip, fmt.Errorf("tripdatastore: could not get Trip: %v", err)
	}
	return trip, nil
}

func (db *TripDatastore) PutTrip(t Trip) error {
	ctx := context.Background()
	k := datastore.NameKey("Trip", t.Event, nil)
	_, err := db.client.Put(ctx, k, &t)
	if err != nil {
		return fmt.Errorf("tripdatastore: could not put Trip: %v", err)
	}
	return nil
}

func (db *TripDatastore) DeleteTrip(event string) error {
	ctx := context.Background()
	k := datastore.NameKey("Trip", event, nil)
	if err := db.client.Delete(ctx, k); err != nil {
		return fmt.Errorf("tripdatastore: could not delete Trip: %v", err)
	}
	return nil
}

func (db *TripDatastore) UpdateTrip(t Trip) error {
	ctx := context.Background()
	k := datastore.NameKey("Trip", t.Event, nil)
	if _, err := db.client.Put(ctx, k, &t); err != nil {
		return fmt.Errorf("tripdatastore: could not update Trip: %v", err)
	}
	return nil
}

func (db *TripDatastore) ListTrips() ([]Trip, error) {
	ctx := context.Background()
	var trips []Trip
	q := datastore.NewQuery("Trip").
		Order("Destination")

	it := db.client.Run(ctx, q)
	for {
		var trip Trip
		_, err := it.Next(&trip)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error fetching next trip: %v", err)
		}
		trips = append(trips, trip)
	}

	return trips, nil
}

func (db *TripDatastore) ListTripsWithDestination(dest string) ([]Trip, error) {
	ctx := context.Background()
	var trips []Trip
	q := datastore.NewQuery("Trip")
	it := db.client.Run(ctx, q)
	for {
		var trip Trip
		_, err := it.Next(&trip)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error fetching next trip: %v", err)
		}
		if strings.Contains(strings.ToLower(trip.Destination), strings.ToLower(dest)) {
			trips = append(trips, trip)
		}
	}

	return trips, nil
}

func (db *TripDatastore) ListTripsCreatedBy(createdBy string) ([]Trip, error) {
	ctx := context.Background()
	var trips []Trip
	q := datastore.NewQuery("Trip").
		Filter("CreatedBy =", createdBy).
		Order("Destination")

	it := db.client.Run(ctx, q)
	for {
		var trip Trip
		_, err := it.Next(&trip)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error fetching next trip: %v", err)
		}
		trips = append(trips, trip)
	}

	return trips, nil
}
