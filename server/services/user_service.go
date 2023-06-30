package services

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/AmirRezaM75/skull-king/contracts"
	"github.com/AmirRezaM75/skull-king/models"
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"html/template"
	"os"
	"strconv"
	"time"
)

type UserService struct {
	userRepository  contracts.UserRepository
	tokenRepository contracts.TokenRepository
}

func NewUserService(userRepository contracts.UserRepository, tokenRepository contracts.TokenRepository) contracts.UserService {
	return UserService{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
	}
}

func (service UserService) Create(email, username, rawPassword string) (*models.User, error) {
	hashedPassword, err := support.HashPassword(rawPassword)

	if err != nil {
		return nil, err
	}

	var user = models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	return service.userRepository.Create(user)
}

func (service UserService) FindByUsername(username string) *models.User {
	return service.userRepository.FindByUsername(username)
}

func (service UserService) FindById(id string) *models.User {
	return service.userRepository.FindById(id)
}

func (service UserService) ExistsByUsername(username string) bool {
	return service.userRepository.ExistsByUsername(username)
}

func (service UserService) ExistsByEmail(email string) bool {
	return service.userRepository.ExistsByEmail(email)
}

func (service UserService) SendEmailVerificationNotification(userId string, email string) error {
	t, err := template.ParseFiles("resources/views/email/email-verification.html")

	if err != nil {
		return errors.New("can't parse HTML file")
	}

	urlGenerator := support.NewUrlGenerator(
		os.Getenv("FRONTEND_URL"),
		os.Getenv("APP_KEY"),
	)

	h := sha1.New()
	h.Write([]byte(email))

	path := fmt.Sprintf(
		"/verify-email/%s/%s",
		userId,
		base64.RawURLEncoding.EncodeToString(h.Sum(nil)),
	)

	verificationURL, _ := urlGenerator.TemporarySignedRoute(
		path,
		map[string]string{},
		time.Now().Add(time.Hour),
	)

	payload := struct {
		AppName string
		Link    string
	}{
		os.Getenv("APP_NAME"),
		verificationURL,
	}

	var body bytes.Buffer

	err = t.Execute(&body, payload)

	if err != nil {
		return err
	}

	m := support.Mail{
		To:      []string{email},
		Subject: "Email Verification",
		Body:    body.String(),
	}

	// TODO: return error
	m.Send()

	return nil
}

func (service UserService) MarkEmailAsVerified(userId string) {
	service.userRepository.UpdateEmailVerifiedAtByUserId(userId, time.Now())
}

func (service UserService) SendResetLink(email string) error {
	templateFile, err := template.ParseFiles("resources/views/email/reset-password-notification.html")

	if err != nil {
		return errors.New("can't parse reset password notification file")
	}

	user := service.userRepository.FindByEmail(email)

	if user == nil {
		return errors.New("we can't find a user with that email address")
	}

	t := service.tokenRepository.FindByEmail(email)

	if t != "" {
		return errors.New("you have exceeded the rate limit. please try again later")
	}

	str := support.Str{}
	token := support.HashHmac(str.Random(40))

	hash, err := support.HashPassword(token)

	if err != nil {
		return err
	}

	lifetime, _ := strconv.Atoi(os.Getenv("RESET_PASSWORD_LINK_LIFETIME"))
	expiration := time.Duration(lifetime) * time.Minute
	err = service.tokenRepository.Create(email, hash, expiration)

	if err != nil {
		return errors.New("couldn't persist token in redis")
	}

	resetPasswordLink := fmt.Sprintf(
		"%s/reset-password?token=%s&email=%s",
		os.Getenv("FRONTEND_URL"),
		token,
		email,
	)

	payload := struct {
		Lifetime int
		Link     string
	}{
		lifetime,
		resetPasswordLink,
	}

	var body bytes.Buffer

	err = templateFile.Execute(&body, payload)

	if err != nil {
		return err
	}

	mail := support.Mail{
		To:      []string{email},
		Subject: "Reset Password Notification",
		Body:    body.String(),
	}

	err = mail.Send()

	if err != nil {
		fmt.Println("SendResetLink: ", err)
		return errors.New("could not send reset password notification")
	}

	return nil
}

func (service UserService) ResetPassword(email, password, token string) error {
	hash := service.tokenRepository.FindByEmail(email)

	if hash == "" {
		return errors.New("the password reset email is invalid")
	}

	err := support.VerifyPassword(hash, token)

	if err != nil {
		fmt.Println(err)
		return errors.New("the password reset token is invalid")
	}

	p, err := support.HashPassword(password)

	if err != nil {
		return errors.New("couldn't not hash password")
	}

	result := service.userRepository.UpdatePasswordByEmail(email, p)

	if result == false {
		return errors.New("couldn't update password")
	}

	_ = service.tokenRepository.DeleteByEmail(email)

	return nil
}
