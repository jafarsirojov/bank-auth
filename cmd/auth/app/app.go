package app

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jafarsirojov/bank-auth/pkg/core/token"
	"github.com/jafarsirojov/bank-auth/pkg/core/users"
	"github.com/jafarsirojov/jwt/pkg/jwt"
	"github.com/jafarsirojov/mux/pkg/mux"
	jwtmidl "github.com/jafarsirojov/mux/pkg/mux/middleware/jwt"
	"github.com/jafarsirojov/rest/pkg/rest"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	router   *mux.ExactMux
	pool     *pgxpool.Pool
	secret   jwt.Secret
	tokenSvc *token.Service
	userSvc  *users.Service
}

func NewServer(router *mux.ExactMux, pool *pgxpool.Pool, secret jwt.Secret, tokenSvc *token.Service, userSvc *users.Service) *Server {
	return &Server{router: router, pool: pool, secret: secret, tokenSvc: tokenSvc, userSvc: userSvc}
}

func (s *Server) Start() {
	s.InitRoutes()
}

type ErrorDTO struct {
	Errors []string `json:"errors"`
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) handleCreateToken(writer http.ResponseWriter, request *http.Request) {
	var body token.RequestDTO
	err := rest.ReadJSONBody(request, &body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		err := rest.WriteJSONBody(writer, &ErrorDTO{
			[]string{"err.json_invalid"},
		})
		log.Print(err)
		return
	}
	response, err := s.tokenSvc.Generate(request.Context(), &body)
	if err != nil {
		switch {
		case err == token.ErrInvalidLogin:
			writer.WriteHeader(http.StatusBadRequest)
			err := rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.login_mismatch"},
			})
			log.Print(err)
		case err == token.ErrInvalidPassword:
			writer.WriteHeader(http.StatusBadRequest)
			err := rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.password_mismatch"},
			})
			log.Print(err)
		default:
			writer.WriteHeader(http.StatusBadRequest)
			err := rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.unknown"},
			})
			log.Print(err)
		}
		return
	}
	err = rest.WriteJSONBody(writer, &response)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleProfile(writer http.ResponseWriter, request *http.Request) {
	authentication, ok := jwtmidl.FromContext(request.Context()).(*token.Payload)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	if authentication.Id == 0 {
		log.Print("admin show users")
		users, err := s.userSvc.Profile(-1)
		if err != nil {
			log.Printf("can't show users: %d", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = rest.WriteJSONBody(writer, users)
		if err != nil {
			log.Print("can't write json body to profile users")
		}
		log.Print("user showing")
		return
	}

	log.Print("admin show profile user")
	users, err := s.userSvc.Profile(authentication.Id)
	if err != nil {
		log.Printf("can't delete user by id: %d", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = rest.WriteJSONBody(writer, users)
	if err != nil {
		log.Print("can't write json body to profile users")
	}
	log.Print("user showing")
	return
}

func (s *Server) handleGetUserByID(writer http.ResponseWriter, request *http.Request) {
	authentication, ok := jwtmidl.FromContext(request.Context()).(*token.Payload)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	//if authentication.Id == 0 {
	//	log.Print("admin show users")
	//	users, err := s.userSvc.Profile(-1)
	//	if err != nil {
	//		log.Printf("can't show users: %d", err)
	//		writer.WriteHeader(http.StatusInternalServerError)
	//		return
	//	}
	//	err = rest.WriteJSONBody(writer, users)
	//	if err != nil {
	//		log.Print("can't write json body to profile users")
	//	}
	//	log.Print("user showing")
	//	return
	//}

	value, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		log.Print("can't delete by id")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		log.Println("Can't strconv.Atoi", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Print("admin show profile user")
	user, err := s.userSvc.GetUserByID(id)
	if err != nil {
		log.Printf("can't delete user by id: %d", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = rest.WriteJSONBody(writer, user)
	if err != nil {
		log.Print("can't write json body to profile users")
	}
	log.Print("user showing")
	return
}

func (s *Server) handlePostSave(writer http.ResponseWriter, request *http.Request) {
	user := users.User{}
	err := rest.ReadJSONBody(request, &user)
	if err != nil {
		log.Print(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.userSvc.Save(user)
	if err != nil {
		log.Printf("can't handle post add user: %d", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (s *Server) handleDeleteUserById(writer http.ResponseWriter, request *http.Request) {
	authentication, ok := jwtmidl.FromContext(request.Context()).(*token.Payload)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	value, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		log.Print("can't delete by id")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		log.Print("can't strconv atoi delete user by id")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if authentication.Id == 0 {
		log.Print("admin deleting user")
		err := s.userSvc.RemoveById(id)
		if err != nil {
			log.Printf("can't delete user by id: %d", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Print("user deleted")
		return
	}
	ownerId := authentication.Id
	if ownerId == id {
		log.Print("user delete himself")
		err := s.userSvc.RemoveById(id)
		if err != nil {
			log.Printf("can't delete user by id: %d", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Print("user deleted")
		return
	}
	writer.WriteHeader(http.StatusBadRequest)
	log.Print("user no deleted")
}
