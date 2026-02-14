package authservice

import (
	"centralService/internal/domain/entities"
	"centralService/internal/domain/enums"
	"centralService/internal/domain/interfaces"
	"centralService/internal/infra/trace"
	"context"
)

type UserService struct {
	storage interfaces.UserStorageInterface
}

func NewUserService(storage interfaces.UserStorageInterface) *UserService {
	return &UserService{
		storage: storage,
	}
}

func (u *UserService) Create(ctx context.Context, user entities.User) (string, error) {
	log := trace.LogWithTraceID("user-service", ctx)

	log.Info("user-service.Create calling")

	user.IsFirstLogin = true

	if err := user.Validate(); err != nil {

		log.Warn("user-service.Create Failed", "error", err.Error())
		err := TranslateUserDomainErrors(err)

		return "", err
	}

	userId, err := u.storage.Create(ctx, user)

	if err != nil {
		err := TranslateStorageErrs(err)

		return "", err
	}

	return userId, nil

}

func (u *UserService) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	log := trace.LogWithTraceID("user-service", ctx)

	log.Info("user-service.GetByEmail calling")

	newUser, err := u.storage.GetByEmail(ctx, email)

	if err != nil {

		err := TranslateStorageErrs(err)

		return nil, err
	}

	return newUser, nil
}

func (u *UserService) GetByID(ctx context.Context, ID string) (*entities.User, error) {
	log := trace.LogWithTraceID("user-service", ctx)

	log.Info("user-service.GetByID calling")

	newUser, err := u.storage.GetByID(ctx, ID)

	if err != nil {

		err := TranslateStorageErrs(err)

		return nil, err
	}

	return newUser, nil
}

func (u *UserService) Update(ctx context.Context, eu *entities.User) error {
	log := trace.LogWithTraceID("user-service", ctx)

	log.Info("user-service.Update calling")

	err := u.storage.Update(ctx, eu)

	if err != nil {

		err := TranslateStorageErrs(err)

		return err
	}

	return nil
}

func (u *UserService) List(ctx context.Context, filter *entities.UserFilter) ([]entities.User, error) {

	log := trace.LogWithTraceID("user-service", ctx).With("method", "List")

	log.InfoContext(ctx, "calling storage")

	users, err := u.storage.List(ctx, filter)

	if err != nil {
		log.ErrorContext(ctx, "failed service", "error", err.Error())
		return nil, TranslateStorageErrs(err)
	}

	return users, nil

}
func (u *UserService) VerifyRole(e *entities.User, target enums.UserPermission) error {

	if err := e.HasPermission(target); err != nil {
		return ErrUserNoHasTargetRole
	}

	return nil
}

func (u *UserService) ConvertPayloadIntoRoles(payload map[string]any) ([]enums.UserPermission, error) {

	v := payload["roles"].([]enums.UserPermission)
	if len(v) == 0 {
		return nil, ErrTokenClaimsInvalid
	}

	return v, nil

}

func (u *UserService) UpdateNameAndPicture(ctx context.Context, id string, name string, picture string) error {

	log := trace.LogWithTraceID("user-service", ctx).With("method", "UpdateNameAndPicture")

	log.Info("calling storage")

	err := u.storage.UpdateNameAndPicture(ctx, id, name, picture)

	if err != nil {
		log.Error("calling storage failed", "error", err.Error())
		return TranslateStorageErrs(err)
	}

	return nil
}
