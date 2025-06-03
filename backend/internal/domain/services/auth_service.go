package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"tracker-backend/internal/domain/dtos"
	"tracker-backend/internal/domain/repository/schemas"
	"tracker-backend/internal/domain/utils"
	"tracker-backend/internal/pkg/logger"
	"tracker-backend/internal/pkg/service"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthorizationService struct {
	userCol *mongo.Collection
	s       SessionProvider
}

var (
	sessionTTL = time.Hour * 24 * 31
)

type PlaylistCreator interface {
	CreateDefault(context.Context, string, dtos.PlaylistCreateRequest) (*schemas.Playlist, error)
}

type SessionProvider interface {
	SetStringTTL(
		context.Context,
		string, string,
		time.Duration,
	) error
	GetString(
		context.Context,
		string,
	) (string, error)
	Invalidate(context.Context, string) error
}

func NewAuthorizationService(userCol *mongo.Collection, s SessionProvider) *AuthorizationService {
	slog.Info("auth service")
	return &AuthorizationService{
		userCol: userCol,
		s:       s,
	}
}

func getSessionIdCacheKey(sessionId string) string {
	return fmt.Sprintf("session:%s", sessionId)
}

func createToken(user *schemas.User, sessionId string) (string, error) {
	claims := utils.JWTClaims{
		UserID:    user.ID,
		SessionID: sessionId,
		Role:      user.Role,
	}
	token, err := utils.CreateTokenFromClaims(claims)
	return token, err
}

func (svc *AuthorizationService) Register(ctx context.Context, req dtos.RegisterRequest) (*schemas.User, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.AuthorizationService.Register"))

	// find similar users
	count, err := svc.userCol.CountDocuments(ctx, bson.M{
		"$or": bson.A{bson.M{"login": req.Login}, bson.M{"email": req.Email}},
	})
	if err != nil {
		_logger.Error("failed count documents", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	// return error if found
	if count > 0 {
		_logger.Warn("same user exists", "request", req)
		return nil, service.ErrExists
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		_logger.Error("error while hashing password", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	// define user
	user := schemas.User{
		ID:              uuid.NewString(),
		Email:           req.Email,
		Login:           req.Login,
		PasswordHash:    passwordHash,
		LikedArtists:    []string{},
		LikedAlbums:     []string{},
		LikedPlaylistId: "nil",
		Role:            schemas.RoleUser,
		CreatedAt:       time.Now(),
	}

	// insert into collection
	_, err = svc.userCol.InsertOne(ctx, user)
	if err != nil {
		slog.Error("failed to insert", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	return &user, nil
}

func (svc *AuthorizationService) Login(
	ctx context.Context,
	req dtos.LoginRequest,
) (string, error) {
	// configure _logger
	_logger := slog.With(slog.String("func", "services.AuthorizationService.Login"))

	// find user
	var user schemas.User
	err := svc.userCol.FindOne(ctx, bson.M{"login": req.Login}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", service.ErrNotFound
		}
		_logger.Error("failed to find user", logger.ErrorAttr(err))
		return "", service.ErrInternal
	}

	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", service.ErrInvalidPassword
	}

	// create token
	sessionId := uuid.NewString()
	token, err := createToken(&user, sessionId)
	if err != nil {
		_logger.Error("failed to create token", logger.ErrorAttr(err))
		return "", service.ErrInternal
	}

	// create session
	key := getSessionIdCacheKey(sessionId)
	err = svc.s.SetStringTTL(ctx, key, token, sessionTTL)
	if err != nil {
		_logger.Warn("failed to set string", logger.ErrorAttr(err))
	}

	return token, nil
}

// Verify returns false if session is invalid
func (svc *AuthorizationService) Verify(ctx context.Context, token string) (*schemas.User, *utils.JWTClaims, bool, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.AuthorizationService.Verify"))

	// parse token
	decoded, claims, err := utils.DecodeToken(token)
	if err != nil {
		_logger.Error("failed to decode token", logger.ErrorAttr(err))
		return nil, nil, false, service.ErrForbidden
	}

	// log claims
	_logger.Debug("decoded claims", "claims", claims, slog.Bool("isValid", decoded.Valid))

	key := getSessionIdCacheKey(claims.SessionID)

	// if token is invalid find session by session id from claims
	cacheToken, _ := svc.s.GetString(ctx, key)
	if cacheToken == "" {
		return nil, nil, false, service.ErrForbidden
	}
	if !decoded.Valid {
		if cacheToken != token {
			return nil, nil, false, service.ErrForbidden
		}
	}

	// get user
	var user schemas.User
	err = svc.userCol.FindOne(ctx, bson.M{"id": claims.UserID}).Decode(&user)
	if err != nil {
		_logger.Error("failed to get user", logger.ErrorAttr(err))
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, false, service.ErrForbidden
		}
		return nil, nil, false, service.ErrInternal
	}

	return &user, claims, true, nil
}

// Refresh returns new token if session exists
func (svc *AuthorizationService) Refresh(
	ctx context.Context,
	token string,
) (string, error) {
	// configure _logger
	_logger := slog.With(slog.String("func", "services.AuthorizationService.Refresh"))

	user, claims, valid, err := svc.Verify(ctx, token)
	if !valid {
		return "", err
	}

	// create new token
	newSessionId := uuid.NewString()
	newToken, err := createToken(user, newSessionId)
	if err != nil {
		_logger.Error("failed to create token", logger.ErrorAttr(err))
		return "", service.ErrInternal
	}

	// update session
	svc.s.Invalidate(ctx, getSessionIdCacheKey(claims.SessionID))
	err = svc.s.SetStringTTL(ctx, getSessionIdCacheKey(newSessionId), newToken, sessionTTL)
	if err != nil {
		_logger.Warn("failed to set string", logger.ErrorAttr(err))
	}

	return newToken, nil
}

// Logout invalidated session
func (svc *AuthorizationService) Logout(
	ctx context.Context,
	token string,
) error {
	// configure _logger
	_logger := slog.With(slog.String("func", "services.AuthorizationService.Refresh"))

	_, claims, err := utils.DecodeToken(token)
	if err != nil {
		_logger.Error("failed to decode token", logger.ErrorAttr(err))
		return service.ErrForbidden
	}

	return svc.s.Invalidate(ctx, getSessionIdCacheKey(claims.SessionID))
}
