package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	authService "tracker-backend/internal/auth"
	auth "tracker-backend/internal/pkg/authorization"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Col       *mongo.Collection
	JwtSecret string
}

func NewUserService(ctx context.Context, db *mongo.Database, jwtSecret string) *UserService {
	// create collection
	col := db.Collection("users")

	// create indices
	err := EnsureIndexes(ctx, col)
	if err != nil {
		panic(err.Error())
	}

	return &UserService{
		Col:       col,
		JwtSecret: jwtSecret,
	}
}

var (
	ErrEmailTaken = errors.New("email already in use")
	ErrLoginTaken = errors.New("login already in use")
	ErrNotFound   = errors.New("user not found")
	ErrForbidden  = errors.New("access denied")
)

func (s *UserService) Register(
	ctx context.Context, credentials RegisterRequest,
) (*User, error) {
	// hash password
	hash, err := auth.HashPassword(credentials.Password)

	if err != nil {
		return nil, fmt.Errorf("error while hashing password")
	}

	// TODO: create user playlist

	// create user schema
	user := &User{
		ID:               uuid.NewString(),
		Login:            credentials.Login,
		Email:            credentials.Email,
		MyChoicePlaylist: "nil",
		PasswordHash:     string(hash),
		CreatedAt:        time.Now(),
		Role:             authService.RoleCustomer,
	}

	// insert user
	_, err = s.Col.InsertOne(ctx, user)

	if err != nil {
		// check constraint error
		if writeErr, ok := err.(mongo.WriteException); ok {
			for _, we := range writeErr.WriteErrors {
				// login constraint error
				if strings.Contains(we.Message, "login") {
					return nil, ErrLoginTaken
				}
				// password constraint error
				if strings.Contains(we.Message, "email") {
					return nil, ErrEmailTaken
				}
			}
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(
	ctx context.Context, credentials LoginRequest,
) (string, error) {
	var user User

	err := s.Col.FindOne(ctx, bson.M{"login": credentials.Login}).Decode(&user)

	if err != nil {
		return "", ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),
		[]byte(credentials.Password)); err != nil {
		return "", ErrForbidden
	}

	return auth.CreateTokenFromUser(
		user.ID, user.Role, user.Email, s.JwtSecret,
	)
}

func (s *UserService) GetByID(
	ctx context.Context, id string,
) (*User, error) {
	var user User
	err := s.Col.FindOne(ctx, bson.M{"id": id}).Decode(&user)

	if err != nil {
		return nil, ErrNotFound
	}

	return &user, nil
}

func (s *UserService) GetAuthDTOByID(
	ctx context.Context, id string, role string,
) (*authService.AuthUser, error) {
	var user authService.AuthUser
	err := s.Col.FindOne(ctx, bson.M{"id": id, "role": role}).Decode(&user)

	if err != nil {
		return nil, ErrNotFound
	}

	return &user, nil
}

func (s *UserService) Update(
	ctx context.Context, id string, req UpdateRequest, allowed bool,
) (*User, error) {
	update := bson.M{}

	if req.Login != nil {
		update["login"] = *req.Login
	}
	if req.Email != nil {
		update["email"] = *req.Email
	}
	if req.Password != nil {
		hashed, err := auth.HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("error while hashing password")
		}
		update["passwordHash"] = hashed
	}
	if req.Role != nil {
		if !allowed {
			return nil, ErrForbidden
		} else {
			update["role"] = *req.Role
		}
	}

	filter := bson.M{"id": id}

	// check uniqueness of login and email
	if login, ok := update["login"]; ok {
		count, err := s.Col.CountDocuments(ctx, bson.M{"login": login, "id": bson.M{"$ne": id}})
		if err != nil {
			return nil, fmt.Errorf("error while checking login uniqueness")
		}
		if count > 0 {
			return nil, ErrLoginTaken
		}
	}

	if email, ok := update["email"]; ok {
		count, err := s.Col.CountDocuments(ctx, bson.M{"email": email, "id": bson.M{"$ne": id}})
		if err != nil {
			return nil, fmt.Errorf("error while checking email uniqueness")
		}
		if count > 0 {
			return nil, ErrEmailTaken
		}
	}

	// update
	_, err := s.Col.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return nil, fmt.Errorf("updating user: %w", err)
	}

	// get updated user
	user, _ := s.GetByID(ctx, id)

	return user, nil
}

func (s *UserService) Delete(
	ctx context.Context, id string,
) error {
	res, err := s.Col.DeleteOne(ctx, bson.M{"id": id})
	if res.DeletedCount < 1 {
		return ErrNotFound
	}
	return err
}
