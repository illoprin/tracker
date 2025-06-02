package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
	authService "tracker-backend/internal/auth"
	"tracker-backend/internal/config"
	auth "tracker-backend/internal/pkg/authorization"
	"tracker-backend/internal/pkg/service"
	playlistType "tracker-backend/internal/playlist/type"
	userType "tracker-backend/internal/user/type"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type PlaylistCreator interface {
	Create(context.Context, playlistType.PlaylistCreateRequest) (*playlistType.Playlist, error)
}

type UserService struct {
	Col *mongo.Collection
	pc  PlaylistCreator
}

func NewUserService(
	ctx context.Context,
	usersCol *mongo.Collection,
	pc PlaylistCreator,
) *UserService {
	return &UserService{
		Col: usersCol,
		pc:  pc,
	}
}

var (
	ErrEmailTaken = errors.New("email already in use")
	ErrLoginTaken = errors.New("login already in use")
)

func (s *UserService) Register(
	ctx context.Context, credentials userType.RegisterRequest,
) (*userType.User, error) {
	// hash password
	hash, err := auth.HashPassword(credentials.Password)

	if err != nil {
		return nil, fmt.Errorf("error while hashing password")
	}

	// create user schema
	user := &userType.User{
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
				// email constraint error
				if strings.Contains(we.Message, "email") {
					return nil, ErrEmailTaken
				}
			}
		}
		return nil, err
	}

	// create default playlist
	defaultPlaylist := playlistType.PlaylistCreateRequest{
		Name:      "Мой выбор",
		UserID:    user.ID,
		IsPublic:  false,
		IsDefault: true,
	}
	p, err := s.pc.Create(ctx, defaultPlaylist)
	if err != nil {
		return user, errors.New("failed to create default playlist")
	}

	// update default playlist pointer
	_, err = s.Col.UpdateOne(ctx,
		bson.M{"id": user.ID},
		bson.M{"$set": bson.M{"myChoicePlaylist": p.ID}},
	)
	if err != nil {
		return user, errors.New("failed to update default playlist")
	}

	return user, nil
}

func (s *UserService) Login(
	ctx context.Context, credentials userType.LoginRequest,
) (string, error) {
	// configure logger
	logger := slog.With(slog.String("function", "track.TrackService.Create"))

	var user userType.User

	err := s.Col.FindOne(ctx, bson.M{"login": credentials.Login}).Decode(&user)

	if err != nil {
		return "", service.ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),
		[]byte(credentials.Password)); err != nil {
		return "", service.ErrAccessDenied
	}

	logger.Info("user authorized", slog.String("id", user.ID))

	return auth.CreateAuthToken(
		user.Role, user.ID, user.Email, os.Getenv(config.JWTSecretEnvName),
	)
}

func (s *UserService) GetByID(
	ctx context.Context, id string,
) (*userType.User, error) {
	var user userType.User
	err := s.Col.FindOne(ctx, bson.M{"id": id}).Decode(&user)

	if err != nil {
		return nil, service.ErrNotFound
	}

	return &user, nil
}

func (s *UserService) GetAuthDTOByID(
	ctx context.Context, id string, role int,
) (*authService.AuthUser, error) {
	var user authService.AuthUser
	err := s.Col.FindOne(ctx, bson.M{"id": id, "role": role}).Decode(&user)

	if err != nil {
		return nil, service.ErrNotFound
	}

	return &user, nil
}

func (s *UserService) Update(
	ctx context.Context, id string, req userType.UpdateRequest, allowed bool,
) (*userType.User, error) {
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
			return nil, service.ErrAccessDenied
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, service.ErrNotFound
		}
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
	// TODO: delete user artists, albums, tracks, clear playlists
	if res.DeletedCount < 1 {
		return service.ErrNotFound
	}
	return err
}
