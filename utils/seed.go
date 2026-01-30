package utils

import (
	"golang-backend/entity"
	"log"
	"strconv"

	"gorm.io/gorm"
)

func SeedRolesAndPermissions(db *gorm.DB) {
	// Define permissions
	permissions := []string{
		"manage_users",
		"manage_roles",
		"view_reports",
	}

	var createdPermissions []*entity.Permission
	for _, permName := range permissions {
		perm := &entity.Permission{Name: permName}
		if err := db.FirstOrCreate(perm, entity.Permission{Name: permName}).Error; err != nil {
			log.Printf("Failed to seed permission %s: %v", permName, err)
		} else {
			createdPermissions = append(createdPermissions, perm)
		}
	}

	// Define Roles
	roles := []string{"admin", "user", "manager"}

	for _, roleName := range roles {
		role := &entity.Role{Name: roleName}
		if err := db.FirstOrCreate(role, entity.Role{Name: roleName}).Error; err != nil {
			log.Printf("Failed to seed role %s: %v", roleName, err)
		}

		// Assign permissions to admin
		if roleName == "admin" {
			if err := db.Model(role).Association("Permissions").Replace(createdPermissions); err != nil {
				log.Printf("Failed to replace permissions for role %s: %v", roleName, err)
			}
		}
	}

	log.Println("✓ Roles and Permissions seeded successfully")
}

func SeedUsers(db *gorm.DB) {
	// Fetch Roles
	var adminRole, userRole, managerRole entity.Role
	db.Where("name = ?", "admin").First(&adminRole)
	db.Where("name = ?", "user").First(&userRole)
	db.Where("name = ?", "manager").First(&managerRole)

	// Pre-hash password "password" to speed up seeding
	hashedPassword, _ := HashPassword("password")

	log.Println("🌱 Seeding users...")

	// 1. Super Admin
	superAdmin := entity.User{
		Name:       "Super Admin",
		Email:      "superadmin@example.com",
		IsVerified: true,
	}
	// Use FirstOrCreate to ensure existence.
	// Note: We need to assign password and roles if creating new.
	// GORM's FirstOrCreate is tricky with related data (Roles).
	// Let's check by email first.
	var existingSuperAdmin entity.User
	if err := db.Where("email = ?", superAdmin.Email).First(&existingSuperAdmin).Error; err != nil {
		// Not found, create
		superAdmin.Password = hashedPassword
		superAdmin.Roles = []*entity.Role{&adminRole, &managerRole}
		if err := db.Create(&superAdmin).Error; err != nil {
			log.Printf("Failed to seed Super Admin: %v", err)
		} else {
			log.Println("✓ Super Admin seeded")
		}
	}

	// 2. Admin
	admin := entity.User{
		Name:       "Admin",
		Email:      "admin@example.com",
		IsVerified: true,
	}
	var existingAdmin entity.User
	if err := db.Where("email = ?", admin.Email).First(&existingAdmin).Error; err != nil {
		admin.Password = hashedPassword
		admin.Roles = []*entity.Role{&adminRole}
		if err := db.Create(&admin).Error; err != nil {
			log.Printf("Failed to seed Admin: %v", err)
		} else {
			log.Println("✓ Admin seeded")
		}
	}

	// 3. Regular Users (Only if count is low)
	var count int64
	db.Model(&entity.User{}).Count(&count)
	if count > 10 {
		return
	}

	for i := 1; i <= 50; i++ {
		email := "user" + strconv.Itoa(i) + "@example.com"
		user := entity.User{
			Name:       "User " + strconv.Itoa(i),
			Email:      email,
			Password:   hashedPassword,
			IsVerified: true,
			Roles:      []*entity.Role{&userRole},
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %s: %v", email, err)
		}
	}
	log.Println("✓ Users seeded successfully")
}
