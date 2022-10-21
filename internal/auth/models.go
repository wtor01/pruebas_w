package auth

import "context"

type User struct {
	ID           string
	Name         string
	Email        string
	AuthID       string
	IsAdmin      bool
	Distributors []Distributor
}

const (
	Admin    = "admin"
	Operator = "operator"
)

var rolePermissionValue = map[string]int{
	Admin:    1,
	Operator: 2,
}

type Role struct {
	ID   string
	Name string
}

type Distributor struct {
	ID   string
	Name string
	Role Role
}

type Repository interface {
	Me(ctx context.Context, id string) (User, error)

	SaveUser(ctx context.Context, user User) (User, error)
	GetDistributor(ctx context.Context, id string) (Distributor, error)
	GetRole(ctx context.Context, id string) (Role, error)
}

type OAuthClient interface {
	CreateUser(ctx context.Context, user User, password string) (string, error)
	DeleteUser(ctx context.Context, id string) error
	VerifyToken(ctx context.Context, token string) (string, error)
}

func (r Role) HaveMorePermission(roleName string) bool {
	preference, ok := rolePermissionValue[roleName]

	if !ok {
		return false
	}

	return r.Preference() <= preference
}

func (r Role) Preference() int {
	preference, _ := rolePermissionValue[r.Name]

	return preference
}

func (u User) HasPermissionInDistributor(distributorID string, roleName string) bool {
	if u.IsAdmin {
		return true
	}

	for _, d := range u.Distributors {
		if d.ID == distributorID && d.Role.HaveMorePermission(roleName) {
			return true
		}
	}

	return false
}

func (u User) GetReadDistributors() []string {
	distributorIds := make([]string, 0)

	for _, d := range u.Distributors {
		if !d.Role.HaveMorePermission(Operator) {
			continue
		}
		distributorIds = append(distributorIds, d.ID)
	}

	return distributorIds
}
