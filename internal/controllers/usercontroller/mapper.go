package usercontroller

import (
	"centralService/internal/controllers"
	"centralService/internal/controllers/usercontroller/models"
	"centralService/internal/domain/entities"
	"centralService/internal/domain/enums"
	"centralService/internal/http/response"
	"fmt"
	"log/slog"
	"net/http"
)

func ParseRoles(r []models.Role) ([]enums.UserPermission, error) {

	roles := map[string]int{
		"admin":       int(enums.USER_PERMISSION_ADMIN),
		"coordinator": int(enums.USER_PERMISSION_COORDINATOR),
		"technical":   int(enums.USER_PERMISSION_TECHNICAL),
	}
	slog.Info("sla o que", "models", r)
	seen := make(map[string]int)
	rolenames := make([]string, 0)
	flags := make([]enums.UserPermission, 0)
	for _, name := range r {
		if _, ok := seen[name.RoleName]; ok {
			return nil, controllers.Wrap(http.StatusBadRequest, &response.ErrorBody{
				Code:    response.ErrBadRequest.Code,
				Message: response.ErrBadRequest.Message,
				Fields: []string{
					fmt.Sprintf("A role %s está duplicada", name.RoleName),
				},
			})
		}

		seen[name.RoleName] = 1

		rolenames = append(rolenames, name.RoleName)

	}
	slog.Info("logando o rolenames", "names", rolenames)

	for _, f := range rolenames {
		if _, ok := roles[f]; !ok {
			return nil, controllers.Wrap(http.StatusBadRequest, &response.ErrorBody{
				Code:    response.ErrBadRequest.Code,
				Message: response.ErrBadRequest.Message,
				Fields: []string{
					fmt.Sprintf("A role %s está invalida", f),
				},
			})
		}

		flags = append(flags, enums.UserPermission(roles[f]))
	}

	return flags, nil
}

func ReparseRoles(eu []enums.UserPermission) ([]string, error) {
	roles := map[enums.UserPermission]string{
		enums.USER_PERMISSION_ADMIN:       "admin",
		enums.USER_PERMISSION_COORDINATOR: "coordinator",
		enums.USER_PERMISSION_TECHNICAL:   "technical",
	}
	rolesString := make([]string, 0, len(eu))

	for _, r := range eu {
		if _, ok := roles[r]; !ok {
			return nil, controllers.Wrap(http.StatusBadRequest, &response.ErrorBody{
				Code:    response.ErrBadRequest.Code,
				Message: response.ErrBadRequest.Message,
				Fields: []string{
					fmt.Sprintf("A role %s está inválida", roles[r]),
				},
			})
		}

		rolesString = append(rolesString, roles[r])
	}

	return rolesString, nil

}

func ModelsUserToEntityUser(mu *models.User) (*entities.User, error) {

	flags, err := ParseRoles(mu.Roles)

	if err != nil {
		return nil, err
	}
	return &entities.User{
		ID:               mu.ID,
		Name:             mu.Name,
		Email:            mu.Email,
		Roles:            flags,
		AssignedToRegion: enums.RemapperRegionFlags[mu.AssignedToRegion],
		Picture:          mu.Picture,
		LastAccessAt:     mu.LastAccessAt,
		CreatedAt:        mu.CreatedAt,
		UpdatedAt:        mu.UpdatedAt,
	}, nil
}

func EntityUserToModelsUser(eu *entities.User) (*models.User, error) {

	roles, err := ReparseRoles(eu.Roles)

	regionFlag := enums.MapperRegionFlags[eu.AssignedToRegion]

	mu := models.User{
		ID:               eu.ID,
		Email:            eu.Email,
		Picture:          eu.Picture,
		AssignedToRegion: regionFlag,
		Name:             eu.Name,
		LastAccessAt:     eu.LastAccessAt,
		CreatedAt:        eu.CreatedAt,
		UpdatedAt:        eu.UpdatedAt,
	}

	if err != nil {
		return nil, err
	}

	for _, r := range roles {
		mu.Roles = append(mu.Roles, models.Role{
			RoleName: r,
		})
	}

	return &mu, nil

}

func ModelUserFilterToEntities(uf *models.UserFilter) (*entities.UserFilter, error) {
	roles, err := ParseRoles(uf.Roles)

	if err != nil {
		return nil, err
	}

	return &entities.UserFilter{
		Name:  uf.Name,
		Email: uf.Email,
		Roles: roles,
		ID:    uf.ID,
	}, nil
}

func EntitiesUsersToModelUsers(eu []entities.User) ([]models.User, error) {

	musers := make([]models.User, 0, len(eu))

	for _, e := range eu {

		mu, err := EntityUserToModelsUser(&e)

		if err != nil {

			return nil, err
		}

		musers = append(musers, *mu)
	}

	return musers, nil
}
