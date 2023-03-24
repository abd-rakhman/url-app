package api

import (
	db "github.com/abd-rakhman/url-app/db/sqlc"
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
	url, err := server.store.CreateNewURLTx(c, db.CreateNewURLRequest{
		HashID: req.HashID,
		URL:    req.Url,
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
	url, err := server.store.GetUrlByHashId(c, req.HashID)
	if err != nil {
		c.JSON(400, gin.H{"error": errorResponse(err)})
		return
	}
	c.JSON(200, gin.H{"url": url.Url})
}
