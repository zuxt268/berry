package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/interface/filter"
	"github.com/zuxt268/berry/internal/repository"
)

type userUsecase struct {
	baseRepository repository.BaseRepository
	userRepository repository.UserRepository
}

type UserUsecase interface {
	GetByUID(ctx context.Context, uid string) (*domain.User, error)
	Gets(ctx context.Context, input domain.GetUsersRequest) (*domain.UsersResponse, error)
	Update(ctx context.Context, input domain.UpdateUserRequest) (*domain.UserResponse, error)
	Create(ctx context.Context, input domain.CreateUserRequest) (*domain.UserResponse, error)
	Delete(ctx context.Context, uid string) error
}

func NewUserUsecase(
	baseRepository repository.BaseRepository,
	userRepository repository.UserRepository,
) UserUsecase {
	return &userUsecase{
		baseRepository: baseRepository,
		userRepository: userRepository,
	}
}

func (u *userUsecase) GetByUID(ctx context.Context, uid string) (*domain.User, error) {
	f := &filter.UserFilter{UID: &uid}
	user, err := u.userRepository.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) Gets(ctx context.Context, input domain.GetUsersRequest) (*domain.UsersResponse, error) {
	f := &filter.UserFilter{
		Name:   input.Name,
		Email:  input.Email,
		Status: input.Status,
	}

	users, err := u.userRepository.List(ctx, f)
	if err != nil {
		return nil, err
	}

	total, err := u.userRepository.Count(ctx, f)
	if err != nil {
		return nil, err
	}

	return toUsersResponse(users, total), nil
}

func (u *userUsecase) Create(ctx context.Context, input domain.CreateUserRequest) (*domain.UserResponse, error) {

	user := &domain.User{
		UID:    uuid.NewString(),
		Name:   input.Name,
		Email:  input.Email,
		Status: domain.UserStatus(input.Status),
	}

	exists, err := u.userRepository.Exists(ctx, &filter.UserFilter{Email: &user.Email})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	created, err := u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return toUserResponse(created), nil
}

func (u *userUsecase) Update(ctx context.Context, input domain.UpdateUserRequest) (*domain.UserResponse, error) {
	f := &filter.UserFilter{UID: &input.UID}

	existing, err := u.userRepository.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Email != nil {
		existing.Email = *input.Email
	}
	if input.Status != nil {
		existing.Status = domain.UserStatus(*input.Status)
	}

	updated, err := u.userRepository.Update(ctx, existing, f)
	if err != nil {
		return nil, err
	}

	return toUserResponse(updated), nil
}

func (u *userUsecase) Delete(ctx context.Context, uid string) error {
	f := &filter.UserFilter{UID: &uid}

	// 存在確認
	if _, err := u.userRepository.Find(ctx, f); err != nil {
		return err
	}

	return u.userRepository.Delete(ctx, f)
}

func toUserResponse(u *domain.User) *domain.UserResponse {
	return &domain.UserResponse{
		ID:        u.ID,
		UID:       u.UID,
		Name:      u.Name,
		Email:     u.Email,
		Status:    int(u.Status),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func toUsersResponse(users []*domain.User, total int64) *domain.UsersResponse {
	userResp := make([]*domain.UserResponse, len(users))
	for i, user := range users {
		userResp[i] = toUserResponse(user)
	}
	return &domain.UsersResponse{
		Users: userResp,
		Paginate: domain.Paginate{
			Total: total,
			Count: int64(len(userResp)),
		},
	}
}
