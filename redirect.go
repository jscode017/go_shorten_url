package main

import (
	//"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Redirect(c *gin.Context) {
	shortenUrl := c.Param("shorten_url")
	log.Println("param: ", shortenUrl)
	if shortenUrl == "" {
		err := errors.New("no shorten url input")
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}

	conn, err := NewRedisConn()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}

	exist, err := CheckKeyExist(conn, shortenUrl)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}
	if exist != 1 {
		err := errors.New("shortenurl not exist or deleted")
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}

	originalUrl, expired, err := GetOriginalUrlAndCheckIfExpire(conn, shortenUrl)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}

	if expired { //do a lazy clean up
		err = DeleteUrlData(conn, shortenUrl)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
			return
		}
		err = errors.New("shorten url expired")
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}
	log.Println("redirecting to: ", originalUrl)
	c.Redirect(http.StatusMovedPermanently, originalUrl)
}
