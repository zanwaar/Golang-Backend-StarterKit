package entity

type Role struct {
	Base
	Name        string        `gorm:"type:varchar(100);uniqueIndex;not null"`
	Permissions []*Permission `gorm:"many2many:role_permissions;"`
}
