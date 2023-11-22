package api

import (
	db "github.com/VrelinnVailati/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccountById)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

func (s *Server) Start(ad string) error {
	return s.router.Run(ad)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
