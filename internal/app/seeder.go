package app

import (
	"context"
	"errors"
	"log"

	"github.com/RizqiSugiarto/coding-test/internal/entity"
	"github.com/RizqiSugiarto/coding-test/internal/repository"
	"github.com/RizqiSugiarto/coding-test/pkg/apperror"
	"golang.org/x/crypto/bcrypt"
)

// SeedUsers seeds the database with sample user data.
func seedUsers(userRepo repository.UserRepo) error {
	ctx := context.Background()

	users := []struct {
		Username string
		Password string
	}{
		{Username: "admin", Password: "admin123"},
		{Username: "user1", Password: "password123"},
		{Username: "user2", Password: "password123"},
		{Username: "testuser", Password: "test123"},
	}

	for _, u := range users {
		// Check if user already exists
		existingUser, err := userRepo.GetByUsername(ctx, u.Username)
		if err != nil && !errors.Is(err, apperror.ErrNotFound) {
			log.Printf("Seeder: error checking user %s: %v", u.Username, err)

			return err
		}

		if existingUser != nil {
			log.Printf("Seeder: user %s already exists, skipping", u.Username)

			continue
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Seeder: error hashing password for user %s: %v", u.Username, err)

			return err
		}

		// Create user
		user := entity.User{
			Username: u.Username,
			Password: string(hashedPassword),
		}

		err = userRepo.Create(ctx, user)
		if err != nil {
			log.Printf("Seeder: error creating user %s: %v", u.Username, err)

			return err
		}

		log.Printf("Seeder: successfully created user %s", u.Username)
	}

	log.Println("Seeder: all users seeded successfully")

	return nil
}
