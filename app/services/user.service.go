package services

import (
	"context"
	"database-ms/app/model"
	"database-ms/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	FindUsersByOrganizationId(context.Context, uuid.UUID) ([]*model.User, error)
	FindByUserEmail(context.Context, string) (*model.User, error)
	FindByUserId(context.Context, uuid.UUID) (*model.User, error)
	Create(context.Context, *model.User) (*mongo.InsertOneResult, error)
	Update(context.Context, *model.User) error
	Delete(context.Context, uuid.UUID) error
	IsUserUnique(context.Context, *model.User) bool
	IsLastAdmin(context.Context, *model.User) (bool, error)
	CreateToken(*gin.Context, *model.User) (string, error)
	HashPassword(string) string
	CheckPasswordHash(string, string) bool
}

type UserService struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewUserService(db *gorm.DB, c *config.Configuration) UserServiceInterface {
	return &UserService{config: c, db: db}
}

func (service *UserService) FindUsersByOrganizationId(ctx context.Context, organizationId uuid.UUID) ([]*model.User, error) {
	var users []*model.User
	result := service.db.Where("organization_id = ?", organizationId).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (service *UserService) FindByUserEmail(ctx context.Context, email string) (*model.User, error) {
	var user *model.User
	result := service.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (service *UserService) FindByUserId(ctx context.Context, userId uuid.UUID) (*model.User, error) {
	user := model.User{}
	user.Id = userId
	result := service.db.First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (service *UserService) Create(ctx context.Context, user *model.User) (*mongo.InsertOneResult, error) {
	result := service.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return nil, nil
}

func (service *UserService) Update(ctx context.Context, user *model.User) error {
	result := service.db.Updates(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (service *UserService) Delete(ctx context.Context, userId uuid.UUID) error {
	user := model.User{}
	user.Id = userId
	result := service.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (service *UserService) IsUserUnique(ctx context.Context, newUser *model.User) bool {
	users, err := service.FindUsersByOrganizationId(ctx, newUser.OrganizationId)
	if err == nil {
		for _, user := range users {
			// Email must be globally unique
			if newUser.Email == user.Email && newUser.Id != user.Id {
				return false
			}
			// Display name must be unique within the organization
			if newUser.DisplayName == user.DisplayName && newUser.OrganizationId == user.OrganizationId && newUser.Id != user.Id {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (service *UserService) IsLastAdmin(ctx context.Context, user *model.User) (bool, error) {
	users, err := service.FindUsersByOrganizationId(ctx, user.OrganizationId)
	if err == nil {
		for _, existingUser := range users {
			if user.Id != existingUser.Id && existingUser.Role == "Admin" {
				return false, nil
			}
		}
		return true, nil
	} else {
		return false, err
	}
}

// ============== Service Helper Method(s) ================

func (service *UserService) CreateToken(c *gin.Context, user *model.User) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["userId"] = user.Id
	atClaims["organizationId"] = user.OrganizationId
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(service.config.AccessSecret))
	if err != nil {
		return "", err
	}
	var expirationDate int = int(time.Now().Add(5 * time.Hour).Unix())
	c.SetCookie("Authorization", token, expirationDate, "/", "", false, true)
	return token, nil
}

func (service *UserService) HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic("Hashing password failed")
	}
	return string(bytes)
}

func (service *UserService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
