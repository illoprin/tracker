package userService

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	hashpassword "tracker-backend/internal/lib/service/hash_password"
	userModel "tracker-backend/internal/user/model"
	userSchema "tracker-backend/internal/user/schema"

	"github.com/golang-jwt/jwt/v5"
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

	// create unique indices
	userSchema.EnsureIndexes(ctx, col)

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
	ctx context.Context, credentials userModel.RegisterRequest,
) (*userSchema.User, error) {
	// hash password
	hash, err := hashpassword.HashPassword(credentials.Password)

	if err != nil {
		return nil, fmt.Errorf("error while hashing password")
	}

	// create user schema
	user := &userSchema.User{
		ID:               uuid.NewString(),
		Login:            credentials.Login,
		Email:            credentials.Email,
		MyChoicePlaylist: "nil",
		PasswordHash:     string(hash),
		CreatedAt:        time.Now(),
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
	ctx context.Context, credentials userModel.LoginRequest,
) (string, error) {
	var user userSchema.User

	err := s.Col.FindOne(ctx, bson.M{"login": credentials.Login}).Decode(&user)
	fmt.Println(err)
	if err != nil {
		return "", ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		return "", ErrForbidden
	}

	fmt.Println(user.ID)

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		// FIX: 24 days to expire token
		"exp": time.Now().Add(24 * time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.JwtSecret))
}

func (s *UserService) GetByID(
	ctx context.Context, id string,
) (*userSchema.User, error) {
	var user userSchema.User
	err := s.Col.FindOne(ctx, bson.M{"id": id}).Decode(&user)

	if err != nil {
		return nil, ErrNotFound
	}

	return &user, nil
}

func (s *UserService) Update(
	ctx context.Context, id string, req userModel.UpdateRequest,
) (*userSchema.User, error) {
	update := bson.M{}

	if req.Login != nil {
		update["login"] = *req.Login
	}
	if req.Email != nil {
		update["email"] = *req.Email
	}
	if req.Password != nil {
		hashed, err := hashpassword.HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("error while hashing password")
		}
		update["passwordHash"] = hashed
	}

	filter := bson.M{"id": id}

	// check uniqueness of login and email
	if login, ok := update["login"]; ok {
		count, err := s.Col.CountDocuments(ctx, bson.M{"login": login, "id": bson.M{"$ne": id}})
		if err != nil {
			return nil, fmt.Errorf("checking login uniqueness: %w", err)
		}
		if count > 0 {
			return nil, ErrLoginTaken
		}
	}

	if email, ok := update["email"]; ok {
		count, err := s.Col.CountDocuments(ctx, bson.M{"email": email, "id": bson.M{"$ne": id}})
		if err != nil {
			return nil, fmt.Errorf("checking email uniqueness: %w", err)
		}
		if count > 0 {
			return nil, ErrEmailTaken
		}
	}

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
