package api

import (
	"fmt"

	db "github.com/abd-rakhman/url-app/db/sqlc"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	Queries *db.Queries
	router  *gin.Engine
}

func NewServer(queries *db.Queries) *Server {
	server := &Server{
		Queries: queries,
	}
	router := gin.Default()

	router.POST("/create", server.createURL)
	router.GET("/:hashID", server.redirectURL)

	server.router = router
	return server
}

func (server *Server) Run(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) string {
	return fmt.Sprint("Error: ", err.Error())
}
