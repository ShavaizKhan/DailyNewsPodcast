package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ShavaizKhan/DailyNewsPodcast/webapp-backend/graph/model"
	"golang.org/x/crypto/bcrypt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Context key for storing user ID
type contextKey string

const userCtxKey contextKey = "user_id"

// NewContextWithUserID creates a new context with the user ID.
func NewContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userCtxKey, userID)
}

// GetUserIDFromContext retrieves the user ID from the context.
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userCtxKey).(string)
	if !ok {
		return "", errors.New("user not authenticated")
	}
	return userID, nil
}

// Resolver root
type Resolver struct {
	Store UserStore
}

// Mutation resolver
func (r *mutationResolver) Signup(ctx context.Context, email string, password string) (string, error) {
	// Hash the password before saving it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %w", err)
	}

	// Create a new user in the database.
	user, err := r.Store.CreateUser(ctx, email, string(hashedPassword), "us", "general")
	if err != nil {
		return "", fmt.Errorf("signup failed: %w", err)
	}

	// Generate and return a JWT for the new user
	return GenerateJWT(user.ID)
}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := r.Store.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate and return a new JWT
	return GenerateJWT(user.ID)
}

func (r *mutationResolver) UpdatePreferences(ctx context.Context, country string, topic string) (*model.Preferences, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err // User not authenticated
	}

	// Update the preferences in the database.
	preferences, err := r.Store.UpdateUserPreferences(ctx, userID, country, topic)
	if err != nil {
		return nil, fmt.Errorf("failed to update preferences: %w", err)
	}

	return preferences, nil
}

// Query resolver
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err // User not authenticated
	}

	// Fetch the user from the database.
	user, err := r.Store.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Note: We're not returning the password hash for security.
	return user, nil
}

func (r *queryResolver) Podcast(ctx context.Context, date *string) (*model.Podcast, error) {
	// Default to today
	if date == nil {
		now := time.Now().Format("2006-01-02")
		date = &now
	}

	// bucket := os.Getenv("S3_BUCKET")
	bucket := "news-podcast-bucket"
	key := fmt.Sprintf("general_podcast_%s.mp3", *date)

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)

	url, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}, s3.WithPresignExpires(15*time.Minute))

	if err != nil {
		return nil, err
	}

	return &model.Podcast{
		Date: *date,
		URL:  url.URL,
	}, nil
}
