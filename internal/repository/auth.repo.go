package repository

import (
	"context"

	"github.com/kitae0522/gommunity/internal/dto"
	"github.com/kitae0522/gommunity/internal/model"
	"github.com/kitae0522/gommunity/pkg/crypt"
	"github.com/kitae0522/gommunity/pkg/utils"
)

type AuthRepository struct {
	client *model.PrismaClient
}

func NewAuthRepository(prismaClient *model.PrismaClient) *AuthRepository {
	return &AuthRepository{client: prismaClient}
}

func (r *AuthRepository) CreateUser(req dto.RegisterRequest) (*model.UsersModel, error) {
	salt := crypt.EncodeBase64(utils.GenerateUUID())
	hashPassword := crypt.NewSHA256(req.Password, salt)

	user, err := r.client.Users.CreateOne(
		model.Users.Handle.Set(req.Handle),
		model.Users.Email.Set(req.Email),
		model.Users.HashPassword.Set(hashPassword),
		model.Users.Salt.Set(salt),
		model.Users.Role.Set(model.UserRolesUser),
		model.Users.Name.Set(req.Name),
	).Exec(context.Background())
	return user, err
}

func (r *AuthRepository) GetUserByEmail(email string) (*model.UsersModel, error) {
	return r.findUserByEmail(email)
}

func (r *AuthRepository) GetUserPasswordByEmail(email string) (*dto.PasswordEntity, error) {
	user, err := r.findUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return &dto.PasswordEntity{
		ID:           user.ID,
		HashPassword: user.HashPassword,
		Salt:         user.Salt,
		Role:         user.Role,
		Handle:       user.Handle,
	}, err
}

func (r *AuthRepository) GetUserPasswordByID(ID string) (*dto.PasswordEntity, error) {
	user, err := r.client.Users.FindUnique(model.Users.ID.Equals(ID)).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return &dto.PasswordEntity{
		ID:           user.ID,
		HashPassword: user.HashPassword,
		Salt:         user.Salt,
		Role:         user.Role,
		Handle:       user.Handle,
	}, err
}

func (r *AuthRepository) UpdateUserHandle(ID string, handle string) error {
	_, err := r.client.Users.FindUnique(
		model.Users.ID.Equals(ID),
	).Update(
		model.Users.Handle.Set(handle),
	).Exec(context.Background())
	return err
}

func (r *AuthRepository) UpdateUserPassword(ID, salt, plainPassword string) error {
	hashPassword := crypt.NewSHA256(plainPassword, salt)
	_, err := r.client.Users.FindUnique(
		model.Users.ID.Equals(ID),
	).Update(
		model.Users.HashPassword.Set(hashPassword),
	).Exec(context.Background())
	return err
}

func (r *AuthRepository) DeleteUser(ID string) (bool, error) {
	_, err := r.client.Users.FindUnique(
		model.Users.ID.Equals(ID),
	).Delete().Exec(context.Background())

	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *AuthRepository) findUserByEmail(email string) (*model.UsersModel, error) {
	user, err := r.client.Users.FindUnique(
		model.Users.Email.Equals(email),
	).Exec(context.Background())
	return user, err
}
