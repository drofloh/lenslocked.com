package models

// import (
// 	"fmt"
// 	"testing"
// 	"time"
// )

// func testingServices() (*Services, error) {

// 	// don't forget to ensure the test DB is created:
// 	// "CREATE DATABASE lenslocked_test;"
// 	const (
// 		host     = "localhost"
// 		port     = 5432
// 		user     = "andrewholford"
// 		password = "blablabla"
// 		dbname   = "lenslocked_test"
// 	)
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname)

// 	services, err := NewServices(psqlInfo)
// 	if err != nil {
// 		return nil, err
// 	}
// 	services.db.LogMode(false)
// 	// Clear user table between tests
// 	services.DestructiveReset()
// 	return services, nil
// }

// func TestCreateUser(t *testing.T) {
// 	services, err := testingServices()
// 	us := services.User
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	user := User{
// 		Name:  "Andrew",
// 		Email: "andrew@email.com",
// 	}
// 	err = us.Create(&user)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if user.ID == 0 {
// 		t.Errorf("Expected ID > 0. Received %d", user.ID)
// 	}
// 	if time.Since(user.CreatedAt) > time.Duration(5*time.Second) {
// 		t.Errorf("Expected CreatedAt to be recent. Received %s", user.CreatedAt)
// 	}
// 	if time.Since(user.UpdatedAt) > time.Duration(5*time.Second) {
// 		t.Errorf("Expected CreatedAt to be recent. Received %s", user.UpdatedAt)
// 	}
// }
