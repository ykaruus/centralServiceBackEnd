package entities

import (
	"centralService/internal/domain"
	"centralService/internal/domain/enums"
	"slices"
	"time"
)

type UserFilter struct {
	ID    string
	Name  string
	Email string
	Roles []enums.UserPermission
}

type User struct {
	ID               string
	Name             string
	Email            string
	Roles            []enums.UserPermission
	AssignedToRegion enums.RegionFlag
	IsFirstLogin     bool
	Picture          string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (u *User) Validate() error {

	return nil
}

func (u *User) HasPermission(target enums.UserPermission) error {

	if !slices.Contains(u.Roles, target) {
		return domain.ErrUserNotHasTargetRole
	}

	return nil
}
