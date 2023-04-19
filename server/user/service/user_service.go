package service

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/AmirRezaM75/skull-king/domain"
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"html/template"
	"os"
	"time"
)

type UserService struct {
	repository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) domain.UserService {
	return UserService{
		repository: userRepository,
	}
}

func (service UserService) Create(email, username, rawPassword string) (*domain.User, error) {
	hashedPassword, err := support.HashPassword(rawPassword)

	if err != nil {
		return nil, err
	}

	var user = domain.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	return service.repository.Create(user)
}

func (service UserService) FindByUsername(username string) *domain.User {
	return service.repository.FindByUsername(username)
}

func (service UserService) ExistsByUsername(username string) bool {
	return service.repository.ExistsByUsername(username)
}

func (service UserService) ExistsByEmail(email string) bool {
	return service.repository.ExistsByEmail(email)
}

func (service UserService) SendEmailVerificationNotification(userId int, email string) error {
	t, err := template.ParseFiles("index.html")

	if err != nil {
		return errors.New("can't parse HTML file")
	}

	urlGenerator := support.UrlGenerator{
		SecretKey: os.Getenv("APP_KEY"),
		BaseURL:   os.Getenv("APP_URL"),
	}

	h := sha1.New()
	h.Write([]byte(email))

	path := fmt.Sprintf(
		"verify-email/%d/%s",
		userId,
		base64.URLEncoding.EncodeToString(h.Sum(nil)),
	)

	verificationURL, _ := urlGenerator.TemporarySignedRoute(
		path,
		map[string]string{"state": "online"}, //TODO: It must be present before signature
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
		Subject: "Register",
		Body:    body.String(),
	}

	m.Send()

	return nil
}
