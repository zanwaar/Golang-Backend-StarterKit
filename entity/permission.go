package entity

type Permission struct {
	Base
	Name  string  `gorm:"type:varchar(100);uniqueIndex;not null"`
	Roles []*Role `gorm:"many2many:role_permissions;"`
}
