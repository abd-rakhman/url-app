package api

import (
	"fmt"

	db "github.com/abd-rakhman/url-app/db/sqlc"
	"github.com/abd-rakhman/url-app/utils"
	"github.com/gin-gonic/gin"
)

func (server *Server) welcome(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Welcome to URL Shortener"})
}

type createURLRequest struct {
	Url    string `json:"url" binding:"required"`
	HashID string `json:"hash_id"`
}

func (server *Server) createURL(c *gin.Context) {
	var req createURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": errorResponse(err)})
		return
	}
	fmt.Printf("req: %+v", req)
	if req.HashID == "" {
		req.HashID = utils.RandomString(6)
	}
	url, err := server.Queries.CreateUrl(c, db.CreateUrlParams{
		Url:    req.Url,
		HashID: req.HashID,
	})
	if err != nil {
		c.JSON(400, gin.H{"error": errorResponse(err)})
		return
	}
	c.JSON(201, gin.H{"hash_id": url.HashID, "url": url.Url})
}

type getURLRequest struct {
	HashID string `uri:"hashID" binding:"required"`
}

func (server *Server) redirectURL(c *gin.Context) {
	var req getURLRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, gin.H{"error": errorResponse(err)})
		return
	}
	url, err := server.Queries.GetUrlByHashId(c, req.HashID)
	if err != nil {
		c.JSON(400, gin.H{"error": errorResponse(err)})
		return
	}
	c.JSON(200, gin.H{"url": url.Url})
}
