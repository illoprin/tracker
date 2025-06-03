package dependencies

import "tracker-backend/internal/domain/services"

type Dependencies struct {
	AuthSvc *services.AuthorizationService
	UserSvc *services.UserService
}
