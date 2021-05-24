package services

import (
	"github.com/zoulongbo/go-mall/models"
	"github.com/zoulongbo/go-mall/repositories"
	"golang.org/x/crypto/bcrypt"
)

type User interface {
	CheckUser(username string, pwd string) (user *models.User, isOk bool)
	AddUser(user *models.User) (id int64, err error)
}

type UserService struct {
	userRepos repositories.User
}

func NewUserService() User {
	return &UserService{userRepos: repositories.NewUserManager()}
}

func (u *UserService) CheckUser(username string, pwd string) (user *models.User, isOk bool) {
	var err error
	user, err = u.userRepos.Select(username)
	if err != nil {
		return
	}
	isOk, err = ValidatePwd(user.HashPassword, pwd)
	if !isOk {
		return &models.User{}, isOk
	}
	return
}

func (u *UserService) AddUser(user *models.User) (id int64, err error) {
	pwdByte, err := GeneratePwd(user.HashPassword)
	if err != nil {
		return
	}
	user.HashPassword = string(pwdByte)
	return u.userRepos.Insert(user)
}

func GeneratePwd(pwd string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
}

func ValidatePwd(hashPwd, pwd string) (isOk bool, err error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(pwd)); err != nil {
		return false, err
	}
	return true, nil
}
