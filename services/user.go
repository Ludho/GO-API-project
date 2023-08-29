package services

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"trades/models"
	"trades/repos"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(context.Context, *models.User) error
	Login(context.Context, *models.User) (string, error)
}

type userServiceImpl struct {
	repo         repos.UserRepository
	bcryptSecret string
	jwtSecret    []byte
}

func NewUserService(repo repos.UserRepository, bcryptSecret string, jwtSecret string) UserService {
	return &userServiceImpl{
		repo:         repo,
		bcryptSecret: bcryptSecret,
		jwtSecret:    []byte(jwtSecret),
	}
}

/*
Create user
*/
func (u *userServiceImpl) Register(c context.Context, user *models.User) error {
	// hashed password
	hashedPassword, err := u.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("error hashing user password: %v", err)
	}
	user.Password = hashedPassword
	return u.repo.CreateUser(c, user)
}

/*
Return user if everything match
*/
func (u *userServiceImpl) Login(c context.Context, user *models.User) (string, error) {

	dbUser, err := u.repo.GetUserByUsername(c, user.Username)

	if err != nil {
		return "", fmt.Errorf("error getting user: %v", err)
	}

	if !CheckPasswordHash(user.Password, dbUser.Password) {
		return "", fmt.Errorf("wrong password for user: %s", user.Username)
	}
	token, err := u.GenerateJWT(user)

	if err != nil {
		return "", fmt.Errorf("error generating user token: %v", err)
	}

	return token, nil

}

/*
Return true if the hashed password and password match
*/
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

/*
Return string of hash password
*/
func (u *userServiceImpl) HashPassword(password string) (string, error) {
	// convert hash cost to int
	key, err := strconv.Atoi(u.bcryptSecret)

	if err != nil {
		return "", err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), key)
	return string(bytes), err
}

/*
	Return string of jwt token
*/

func (u *userServiceImpl) GenerateJWT(user *models.User) (string, error) {
	// setup jwt token
	claims := jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// encode secret
	signedToken, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("error generating JWT: %w", err)
	}

	return signedToken, nil
}
