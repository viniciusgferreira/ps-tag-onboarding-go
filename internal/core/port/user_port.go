package port

import (
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/model"
)

// remove interface
type UserService interface {
	Find(ctx *gin.Context, id string) (*model.User, error)
	Save(ctx *gin.Context, u model.User) (*model.User, error)
	Update(ctx *gin.Context, u model.User) (*model.User, error)
}

// move to service
type UserRepository interface {
	FindById(ctx *gin.Context, id string) (*model.User, error)
	Save(ctx *gin.Context, u model.User) (*model.User, error)
	Update(ctx *gin.Context, u model.User) (*model.User, error)
	ExistsByFirstNameAndLastName(ctx *gin.Context, u model.User) bool
}
