package app

import (
	"github.com/jafarsirojov/bank-auth/pkg/core/token"
	"github.com/jafarsirojov/mux/pkg/mux/middleware/authenticated"
	"github.com/jafarsirojov/mux/pkg/mux/middleware/jwt"
	"github.com/jafarsirojov/mux/pkg/mux/middleware/logger"
	"reflect"
)

func (s *Server) InitRoutes() {
	var jwtMiddleware = jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret)
	s.router.POST(
		"/api/tokens",
		s.handleCreateToken,
		logger.Logger("TOKEN"),
	)
	s.router.GET(
		"/api/users",
		s.handleProfile,
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwtMiddleware,
		logger.Logger("USERS"),
	)
	s.router.POST(
		"/api/users",
		s.handlePostSave,
		logger.Logger("POST USERS"),
	)

	s.router.DELETE(
		"/api/users/{id}",
		s.handleDeleteUserById,
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwtMiddleware,
		logger.Logger("DELETE"),
	)
}


