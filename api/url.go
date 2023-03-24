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
	Url       string `json:"url" binding:"required"`
	HashID    string `json:"hash_id"`
	ExpiresAt int64  `json:"expires_at"`
}

func (server *Server) createURL(c *gin.Context) {
	var req createURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": errorResponse(err)})
		return
	}
	fmt.Println(req)
	url, err := server.store.CreateNewURLTx(c, db.CreateNewURLRequest{
		HashID:    req.HashID,
		URL:       req.Url,
		ExpiresAt: req.ExpiresAt,
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

type cleanExpiredURLRequest struct {
	Secret string `json:"secret"`
}

func (server *Server) cleanExpiredURLs(c *gin.Context) {
	var req cleanExpiredURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": errorResponse(err)})
		return
	}
	config, err := utils.LoadConfig("..")
	if err != nil {
		c.JSON(500, gin.H{"error": errorResponse(err)})
		return
	}
	if req.Secret != config.CleanSecretKey {
		c.JSON(400, gin.H{"error": "Incorrect secret key"})
		return
	}
	err = server.store.DeleteExpiredUrls(c)
	if err != nil {
		c.JSON(400, gin.H{"error": errorResponse(err)})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
