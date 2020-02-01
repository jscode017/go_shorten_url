package main

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func DealWithAddUrlRequest(c *gin.Context) { //please use an url that starts with http:// or https://
	log.Println("receive add url request")
	var requestBody map[string]interface{}
	err := json.NewDecoder(c.Request.Body).Decode(&requestBody)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}
	var urlData UrlData
	if url, ok := requestBody["url"]; ok {
		urlData.OriginalUrl = url.(string)
		urlData.ExpireTime = "5000" //TODO: might change it later
	} else {
		err = errors.New("No original url input")
		c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}

	shortenUrl, err := AddUrl(urlData)
	if err != nil {
		log.Print(err)
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"result": true, "shortenurl": shortenUrl})
	return
}

func DealWithDeleteUrlRequest(c *gin.Context) {
	log.Println("receive delete url request")
	var requestBody map[string]interface{}
	err := json.NewDecoder(c.Request.Body).Decode(&requestBody)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}
	var shortenUrl string
	if url, ok := requestBody["url"]; ok {
		shortenUrl = url.(string)
	} else {
		err = errors.New("No shorten url input")
		c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}

	err = RemoveUrl(shortenUrl)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"result": true, "shortenurl": shortenUrl + " deleted"})
	return

}
func AddUrl(urlData UrlData) (string, error) {
	shortenUrl := GenerateKey(urlData.OriginalUrl)
	conn, err := NewRedisConn()
	if err != nil {
		return "", err
	}

	exist, err := CheckKeyExist(conn, shortenUrl)
	if err != nil {
		return "", err
	}

	for exist == 1 {
		randomNum := RandonIntFromTime()
		newShortenUrl := GenerateKey(urlData.OriginalUrl + strconv.Itoa(randomNum)) //keep trying to get an unique shortenUrl
		exist, err = CheckKeyExist(conn, shortenUrl)
		if err != nil {
			return "", err
		}

		if exist != 1 {
			shortenUrl = newShortenUrl
		}
	}

	urlData.ShortenUrl = shortenUrl
	urlData.CreationTime = ConvertTimeToStr(time.Now())

	err = AddUrlToDb(conn, urlData) //may still have concurrency problem for conflict shortenUrls, TODO: use a key generation service
	if err != nil {
		return "", err
	}
	return shortenUrl, nil
}

func RemoveUrl(shortenUrl string) error {
	conn, err := NewRedisConn()
	if err != nil {
		return err
	}

	err = DeleteUrlData(conn, shortenUrl)
	if err != nil {
		return err
	}

	return nil
}
