package dependencies

import (
	"tracker-backend/internal/domain/repository"
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/infrastructure/redis"
)

func InitDependencies(repo *repository.Repository, redisClient *redis.RedisClient) *Dependencies {
	authSvc := services.NewAuthorizationService(repo.UsersCol, redisClient)
	userSvc := services.NewUserService(repo.UsersCol)

	return &Dependencies{
		AuthSvc: authSvc,
		UserSvc: userSvc,
	}
}
