package services

import (
	"context"
	"errors"
	"log/slog"
	"tracker-backend/internal/domain/dtos"
	"tracker-backend/internal/domain/repository/schemas"
	"tracker-backend/internal/domain/utils"
	"tracker-backend/internal/pkg/logger"
	"tracker-backend/internal/pkg/service"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserFlusher interface {
	FlushUserData(ctx context.Context, id string) error
}

type UserService struct {
	userCol *mongo.Collection
	uf      UserFlusher
}

func NewUserService(userCol *mongo.Collection, userFlusher UserFlusher) *UserService {
	return &UserService{
		userCol: userCol,
		uf:      userFlusher,
	}
}

func (svc *UserService) GetByID(ctx context.Context, userId string) (*schemas.User, error) {
	_logger := slog.With(slog.String("func", "services.UserService.GetByID"))

	var user schemas.User
	err := svc.userCol.FindOne(ctx, bson.M{"id": userId}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, service.ErrNotFound
		}
		_logger.Error("failed to find user", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	return &user, nil
}

func (svc *UserService) UpdateByID(
	ctx context.Context,
	userId string,
	req dtos.UserUpdateRequest,
	roleChangingAllowed bool,
) (*schemas.User, error) {
	// configure logger
	_logger := slog.With(slog.String("func", "services.UserService.UpdateByID"))

	filter := bson.M{"id": userId}
	updates := bson.M{}

	// define updates
	if req.Login != nil {
		updates["login"] = *req.Login
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Password != nil {
		p, err := utils.HashPassword(*req.Password)
		if err != nil {
			_logger.Error("failed to hash password", logger.ErrorAttr(err))
			return nil, service.ErrInternal
		}
		updates["passwordHash"] = p
	}
	// HACK: to change role you need 'Allow-Access' header
	if req.Role != nil {
		if !roleChangingAllowed {
			return nil, service.ErrForbidden
		}
		updates["role"] = *req.Role
	}

	var updated schemas.User
	err := svc.userCol.FindOneAndUpdate(
		ctx,
		filter,
		bson.M{"$set": updates},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, service.ErrNotFound
		}
		_logger.Error("failed to update user", logger.ErrorAttr(err))
		return nil, service.ErrInternal
	}

	return &updated, nil
}

func (svc *UserService) DeleteByID(ctx context.Context, userId string) error {
	// configure logger
	_logger := slog.With(slog.String("func", "services.UserService.DeleteByID"))

	// flush related data
	if err := svc.uf.FlushUserData(ctx, userId); err != nil {
		return service.ErrInternal
	}

	// delete user document
	res, err := svc.userCol.DeleteOne(ctx, bson.M{"id": userId})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return service.ErrNotFound
		}
		_logger.Error("failed to update user", logger.ErrorAttr(err))
		return service.ErrInternal
	}

	if res.DeletedCount == 0 {
		return service.ErrNotFound
	}

	return nil
}
