package session

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"grader/pkg/server/user"
	"time"
)

type JWT struct {
	client    *redis.Client
	jwtSecret []byte
}

type JWTClaims struct {
	ID   string       `json:"sid"`
	User *user.Claims `json:"user"`
	jwt.StandardClaims
}

func NewSessionJWT(jwtSecret string, client *redis.Client) *JWT {
	return &JWT{
		jwtSecret: []byte(jwtSecret),
		client:    client,
	}
}

func (sm *JWT) parseSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)

	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}

	return sm.jwtSecret, nil
}

func (sm *JWT) Check(token string) (*Session, error) {

	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	payload := &JWTClaims{}

	tkn, err := jwt.ParseWithClaims(token, payload, sm.parseSecretGetter)

	if err != nil {
		return nil, fmt.Errorf("cant parse json %v", err)
	}

	if !tkn.Valid {
		return nil, fmt.Errorf("invalid jwt token %v", err)
	}

	key := fmt.Sprintf("%s", payload.ID)

	userId, err := sm.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, errors.New("session expired")
	}

	u := payload.User

	if u.ID != userId {
		return nil, err
	}

	return &Session{
		ID:   payload.ID,
		User: u,
	}, nil
}

func (sm *JWT) Create(u *user.User) (string, error) {
	sessId := uuid.New().String()

	err := sm.client.Set(context.Background(), sessId, u.ID, 24*time.Hour).Err()
	if err != nil && err != redis.Nil {
		return "", err
	}

	data := &JWTClaims{
		ID: sessId,
		User: &user.Claims{
			ID:       u.ID,
			Username: u.Username,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	signedToken, err := token.SignedString(sm.jwtSecret)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (sm *JWT) Delete(token string) error {
	payload := &JWTClaims{}

	_, err := jwt.ParseWithClaims(token, payload, sm.parseSecretGetter)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s", payload.ID)

	err = sm.client.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}

	return nil
}

func (sm *JWT) Get(sess *Session) (string, error) {
	sessId, err := sm.client.Get(context.Background(), sess.ID).Result()
	if err != nil {
		return "", err
	}

	return sessId, nil
}

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(Key).(*Session)

	if !ok || sess == nil {
		return nil, ErrNoAuth
	}

	return sess, nil
}
