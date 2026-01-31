package entity_test

import (
	"golang-backend/entity"
	"testing"
)

func TestUser_HasRole(t *testing.T) {
	adminRole := &entity.Role{Name: "admin"}
	userRole := &entity.Role{Name: "user"}

	user := &entity.User{
		Roles: []*entity.Role{adminRole},
	}

	if !user.HasRole("admin") {
		t.Error("User should have admin role")
	}

	if user.HasRole("manager") {
		t.Error("User should not have manager role")
	}

	user.AssignRole(userRole)
	if !user.HasRole("user") {
		t.Error("User should have user role after assignment")
	}
}

func TestUser_HasPermission(t *testing.T) {
	perm1 := &entity.Permission{Name: "edit_post"}
	perm2 := &entity.Permission{Name: "delete_post"}

	role := &entity.Role{
		Name:        "editor",
		Permissions: []*entity.Permission{perm1},
	}

	user := &entity.User{
		Roles: []*entity.Role{role},
	}

	if !user.HasPermission("edit_post") {
		t.Error("User should have edit_post permission")
	}

	if user.HasPermission("delete_post") {
		t.Error("User should not have delete_post permission")
	}

	role.Permissions = append(role.Permissions, perm2)
	if !user.HasPermission("delete_post") {
		t.Error("User should have delete_post permission after updating role")
	}
}
