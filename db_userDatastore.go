package ridemerge

import (
	"fmt"
	"log"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/datastore"

	"golang.org/x/net/context"
)

type UserDatastore struct {
	client *datastore.Client
}

var _ UserDatabase = &UserDatastore{}

func NewUserDatastore(client *datastore.Client) (UserDatabase, error) {
	ctx := context.Background()
	t, err := client.NewTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("userdatastore: could not connect: %v", err)
	}
	if err := t.Rollback(); err != nil {
		return nil, fmt.Errorf("userdatastore: could not connect: %v", err)
	}
	return &UserDatastore{
		client: client,
	}, nil
}

// Close closes the database.
func (db *UserDatastore) Close() {
	// No op.
}

func (db *UserDatastore) datastoreKey(email string) *datastore.Key {
	return datastore.NameKey("User", email, nil)
}

func (db *UserDatastore) GetUser(email string) (User, error) {
	ctx := context.Background()
	k := datastore.NameKey("User", email, nil)
	var user User
	if err := db.client.Get(ctx, k, &user); err != nil {
		return user, fmt.Errorf("userdatastore: could not get User: %v", err)
	}
	return user, nil
}

func (db *UserDatastore) PutUser(u User) error {
	ctx := context.Background()
	k := datastore.NameKey("User", u.Email, nil)
	_, err := db.client.Put(ctx, k, &u)
	if err != nil {
		return fmt.Errorf("userdatastore: could not put User: %v", err)
	}
	return nil
}

func (db *UserDatastore) DeleteUser(email string) error {
	ctx := context.Background()
	k := datastore.NameKey("User", email, nil)
	if err := db.client.Delete(ctx, k); err != nil {
		return fmt.Errorf("userdatastore: could not delete User: %v", err)
	}
	return nil
}

func (db *UserDatastore) UpdateUser(u User) error {
	ctx := context.Background()
	k := datastore.NameKey("User", u.Email, nil)
	if _, err := db.client.Put(ctx, k, &u); err != nil {
		return fmt.Errorf("userdatastore: could not update User: %v", err)
	}
	return nil
}

func (db *UserDatastore) ListUsers() ([]User, error) {
	ctx := context.Background()
	users := make([]User, 0)
	q := datastore.NewQuery("User").
		Order("LastName")

	it := db.client.Run(ctx, q)
	for {
		var user User
		_, err := it.Next(&user)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error fetching next user: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}
