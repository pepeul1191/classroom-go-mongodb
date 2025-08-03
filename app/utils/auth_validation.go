// utils/validation.go
package utils

import (
	"classroom/app/models"
	"errors"
)

func ValidateRoles(roles []models.Role) error {
	if len(roles) == 0 {
		return errors.New("el usuario debe tener al menos un rol asignado")
	}

	for _, role := range roles {
		if len(role.Permissions) == 0 {
			return errors.New("cada rol debe tener al menos un permiso")
		}

		for _, permission := range role.Permissions {
			if permission == "" {
				return errors.New("los permisos no pueden estar vacíos")
			}
		}
	}

	return nil
}
