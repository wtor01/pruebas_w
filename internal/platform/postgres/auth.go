package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/auth"
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ModelEntity
	AuthID       string         `gorm:"column:auth_id"`
	Name         string         `gorm:"column:name"`
	Email        string         `gorm:"column:email"`
	IsAdmin      bool           `gorm:"column:is_admin"`
	Roles        []*Role        `gorm:"many2many:user_distributor_roles;"`
	Distributors []*Distributor `gorm:"many2many:user_distributor_roles;"`
}

func (u User) toDomain() auth.User {
	permissions := make([]auth.Distributor, 0, cap(u.Distributors))
	for _, distributor := range u.Distributors {
		permission := auth.Distributor{
			ID:   distributor.ID.String(),
			Name: distributor.Name,
			Role: auth.Role{},
		}
		for _, role := range distributor.Roles {
			permission.Role = auth.Role{
				ID:   role.ID.String(),
				Name: role.Name,
			}
			permissions = append(permissions, permission)

		}
	}
	return auth.User{
		ID:           u.ID.String(),
		Name:         u.Name,
		Email:        u.Email,
		AuthID:       u.AuthID,
		IsAdmin:      u.IsAdmin,
		Distributors: permissions,
	}
}

func (d Distributor) toAuthDomain() auth.Distributor {
	return auth.Distributor{
		ID:   d.ID.String(),
		Name: d.Name,
	}
}

type Role struct {
	ModelEntity
	Name         string         `gorm:"column:name"`
	Users        []*User        `gorm:"many2many:user_distributor_roles;"`
	Distributors []*Distributor `gorm:"many2many:user_distributor_roles;"`
}

func (r Role) toDomain() auth.Role {
	return auth.Role{
		ID:   r.ID.String(),
		Name: r.Name,
	}
}

type UserDistributorRole struct {
	RoleID        string      `gorm:"primaryKey"`
	UserID        string      `gorm:"primaryKey"`
	DistributorID string      `gorm:"primaryKey"`
	Role          Role        `gorm:"foreignKey:RoleID"`
	User          User        `gorm:"foreignKey:UserID"`
	Distributor   Distributor `gorm:"foreignKey:DistributorID"`
}

func (UserDistributorRole) TableName() string {
	return "user_distributor_roles"
}

type AuthRepository struct {
	db        *gorm.DB
	dbTimeout time.Duration
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db, dbTimeout: time.Second * 5}
}

func (repository AuthRepository) SaveUser(ctx context.Context, user auth.User) (auth.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, repository.dbTimeout)
	defer cancel()

	err := repository.db.WithContext(ctxTimeout).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&UserDistributorRole{}, "user_id = ?", user.ID).Error; err != nil {
			return err
		}
		if err := tx.Create(&User{
			ModelEntity: ModelEntity{
				ID:        uuid.FromStringOrNil(user.ID),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			AuthID:  user.AuthID,
			Name:    user.Name,
			Email:   user.Email,
			IsAdmin: user.IsAdmin,
		}).Error; err != nil {
			return err
		}

		if len(user.Distributors) == 0 {
			return nil
		}

		params := make([]interface{}, 0, cap(user.Distributors)*3)

		query := fmt.Sprintf("INSERT INTO %s (user_id, distributor_id, role_id) VALUES ", (UserDistributorRole{}).TableName())

		for _, p := range user.Distributors {
			query += fmt.Sprintf("(?, ?, ?),")
			params = append(params, user.ID, p.ID, p.Role.ID)
		}
		query = query[:len(query)-1]

		if err := tx.Raw(query, params...).Scan(&map[string]interface{}{}).Error; err != nil {
			return err
		}

		return nil
	})

	return user, err
}

func (repository AuthRepository) GetDistributor(ctx context.Context, id string) (auth.Distributor, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, repository.dbTimeout)
	defer cancel()

	var distributor Distributor
	result := repository.db.WithContext(ctxTimeout).First(&distributor, "id = ?", id)

	if result.Error != nil {
		return auth.Distributor{}, result.Error
	}

	return distributor.toAuthDomain(), nil
}

func (repository AuthRepository) GetRole(ctx context.Context, id string) (auth.Role, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, repository.dbTimeout)
	defer cancel()

	var role Role
	result := repository.db.WithContext(ctxTimeout).First(&role, "id = ?", id)

	if result.Error != nil {
		return auth.Role{}, result.Error
	}

	return role.toDomain(), nil
}

func (repository AuthRepository) Me(ctx context.Context, id string) (auth.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, repository.dbTimeout)
	defer cancel()

	var user User
	result := repository.db.WithContext(ctxTimeout).Preload("Distributors.Roles").Preload("Distributors").First(&user, "id = ? ", id)

	if result.Error != nil {
		return auth.User{}, result.Error
	}

	return user.toDomain(), nil
}
