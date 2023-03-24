package api

import (
	"fmt"

	db "github.com/abd-rakhman/url-app/db/sqlc"
	"github.com/abd-rakhman/url-app/utils"
	"github.com/gin-gonic/gin"
)

type createURLRequest struct {
	URL string `json:"url" binding:"required"`
}

func (server *Server) createURL(c *gin.Context) {
	var req createURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": errorResponse(err)})
		return
	}
	randomHashID := utils.RandomString(6)
	fmt.Println(randomHashID)
	url, err := server.Queries.CreateUrl(c, db.CreateUrlParams{
		HashID: randomHashID,
		Url:    req.URL,
	})
	if err != nil {
		c.JSON(400, gin.H{"error": errorResponse(err)})
		return
	}
	c.JSON(200, gin.H{"hash_id": url.HashID, "url": url.Url})
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
