package api

import (
	"fmt"

	db "github.com/abd-rakhman/url-app/db/sqlc"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.GET("/welcome", server.welcome)
	router.POST("/create", server.createURL)
	router.GET("/:hashID", server.redirectURL)
	router.POST("/clean-expired", server.cleanExpiredURLs)

	server.router = router
	return server
}

func (server *Server) Run(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) string {
	return fmt.Sprint("Error: ", err.Error())
}
