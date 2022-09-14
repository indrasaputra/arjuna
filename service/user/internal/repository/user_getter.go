package repository

import (
	"context"

	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/service"
)

// UserGetter is responsible to connect user with repositories.
type UserGetter struct {
	keycloak service.GetUserRepository
	postgres service.GetUserRepository
}

// NewUserGetter creates an instance of UserGetter.
func NewUserGetter(keycloak, postgres service.GetUserRepository) *UserGetter {
	return &UserGetter{
		keycloak: keycloak,
		postgres: postgres,
	}
}

// GetAll gets all users from Keycloak and Postgres.
func (ug *UserGetter) GetAll(ctx context.Context) ([]*entity.User, error) {
	kcUsers, err := ug.keycloak.GetAll(ctx)
	if err != nil {
		return []*entity.User{}, err
	}
	pgUsers, err := ug.postgres.GetAll(ctx)
	if err != nil {
		return []*entity.User{}, err
	}
	return combineUsersKcAndPg(kcUsers, pgUsers), nil
}

func combineUsersKcAndPg(kcUsers, pgUsers []*entity.User) []*entity.User {
	mapping := make(map[string]*entity.User)
	for _, u := range kcUsers {
		mapping[u.Email] = u
	}

	for i, u := range pgUsers {
		pgUsers[i].KCID = mapping[u.Email].KCID
	}
	return pgUsers
}
