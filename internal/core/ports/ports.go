package ports

import (
	"github.com/gin-gonic/gin"
	"github.com/viniciusgferreira/ps-tag-onboarding-go/internal/core/domain/models"
)

type UserService interface {
	Find(ctx *gin.Context, id string) (*models.User, error)
	Save(ctx *gin.Context, u models.User) (*models.User, error)
	Update(ctx *gin.Context, u models.User) (*models.User, error)
}

type UserRepository interface {
	FindById(ctx *gin.Context, id string) (*models.User, error)
	Save(ctx *gin.Context, u models.User) (*models.User, error)
	Update(ctx *gin.Context, u models.User) (*models.User, error)
	ExistsByFirstNameAndLastName(ctx *gin.Context, firstName, lastName string) (bool, error)
}
