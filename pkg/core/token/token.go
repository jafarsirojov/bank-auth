package token

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jafarsirojov/jwt/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type Service struct {
	secret []byte
	pool   *pgxpool.Pool
}

func NewService(secret []byte, pool *pgxpool.Pool) *Service {
	return &Service{secret: secret, pool: pool}
}

type Payload struct {
	Id    int `json:"id"`
	Exp   int64 `json:"exp"`
	Phone int   `json:"phone"`
}

type RequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseDTO struct {
	Token string `json:"token"`
}

var ErrInvalidLogin = errors.New("invalid password")
var ErrInvalidPassword = errors.New("invalid password")



func (s *Service) Generate(context context.Context, request *RequestDTO) (response ResponseDTO, err error) {
	var pass string
	var id int
	var phone int
	log.Print("a")
	log.Print(request.Username, "log")
	log.Print(request.Password, "pass")
	log.Print(context)
	err = s.pool.QueryRow(context,
		`SELECT password, id, phone 
		FROM users
		WHERE removed = FALSE AND login = $1`,
		request.Username,
	).Scan(
		&pass,
		&id,
		&phone,
	)
	log.Print("queryrow")
	log.Print("b")
	if err != nil {
		err = ErrInvalidLogin
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	log.Print(string(hash))
log.Print("c")
	err = bcrypt.CompareHashAndPassword(hash, []byte(request.Password))
	log.Print("d")
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		err = ErrInvalidPassword
		return
	}

	response.Token, err = jwt.Encode(Payload{
		Id:    id,
		Exp:   time.Now().Add(time.Hour).Unix(),
		Phone: phone,
	}, s.secret)
	log.Print("e finish")
	return
}
