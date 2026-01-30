package entity

import (
	"time"
)

type User struct {
	Base
	Name             string `gorm:"type:varchar(100);not null"`
	Email            string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password         string `gorm:"not null"`
	IsVerified       bool   `gorm:"default:false"`
	VerificationCode string `gorm:"type:varchar(6)"`
	ResetToken       string `gorm:"type:varchar(6)"`
	ResetTokenExpiry time.Time
	Roles            []*Role `gorm:"many2many:user_roles;"`
}

// Helper methods for Role & Permission checks

func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

func (u *User) HasPermission(permissionName string) bool {
	for _, role := range u.Roles {
		for _, permission := range role.Permissions {
			if permission.Name == permissionName {
				return true
			}
		}
	}
	return false
}

func (u *User) AssignRole(role *Role) {
	u.Roles = append(u.Roles, role)
}
