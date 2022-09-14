package service

// import (
// 	"context"
// 	"log"
// )

// // DeleteUser defines the interface to delete a user.
// type DeleteUser interface {
// 	// HardDeleteAll hard-deletes all users.
// 	HardDeleteAll(ctx context.Context, key string) error
// }

// // DeleteUserRepository defines the interface to delete user from the repository.
// type DeleteUserRepository interface {
// 	// GetAll gets all users available in repository.
// 	// If there isn't any user in repository, it returns empty list of user and nil error.
// 	GetAll(ctx context.Context) ([]*entity.User, error)
// 	// DeleteByKey deletes a single user from the repository.
// 	// If the user can't be found, it doesn't return error.
// 	DeleteByKey(ctx context.Context, key string) error
// }

// // UserDeleter is responsible for deleting a user.
// type UserDeleter struct {
// 	repo      DeleteUserRepository
// 	publisher UserPublisher
// }

// // NewUserDeleter creates an instance of UserDeleter.
// func NewUserDeleter(repo DeleteUserRepository, publisher UserPublisher) *UserDeleter {
// 	return &UserDeleter{
// 		repo:      repo,
// 		publisher: publisher,
// 	}
// }

// // DeleteByKey deletes a user by its key.
// // It only deletes disabled user.
// func (td *UserDeleter) DeleteByKey(ctx context.Context, key string) error {
// 	user, err := td.repo.GetByKey(ctx, key)
// 	if err != nil {
// 		return err
// 	}
// 	if user.IsEnabled {
// 		return entity.ErrProhibitedToDelete()
// 	}
// 	if err := td.repo.DeleteByKey(ctx, key); err != nil {
// 		return err
// 	}
// 	if err := td.publisher.Publish(ctx, entity.EventUserDeleted(user)); err != nil {
// 		log.Printf("publish on user deleter error: %v", err)
// 	}
// 	return nil
// }
