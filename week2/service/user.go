package service

import (
	"awesomeProject/webook/internal/domain"
	"awesomeProject/webook/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")
var ErrInvalidUserNotFund = errors.New("不存在该用户")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	//先找用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	//你要考虑加密放在哪些的问题
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
	//println(1111)
	//svc.repo.Create(ctx, u)

	//return svc.repo.Create(ctx, u)
}

func (svg *UserService) EditUserProfile(ctx context.Context, id int64, u domain.User) error {
	return svg.repo.Edit(ctx, id, u)
}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	//先找用户
	u, err := svc.repo.FindById(ctx, id)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserNotFund
	}
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
