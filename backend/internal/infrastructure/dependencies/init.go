package dependencies

import (
	"tracker-backend/internal/domain/repository"
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/infrastructure/redis"
)

func InitDependencies(repo *repository.Repository, redisClient *redis.RedisClient) *Dependencies {
	authSvc := services.NewAuthorizationService(repo.UsersCol, redisClient)

	return &Dependencies{
		AuthSvc: authSvc,
	}
}
