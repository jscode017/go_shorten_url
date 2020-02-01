package main

import (
	"github.com/gomodule/redigo/redis"
	//"fmt"
	"errors"
	//"strings"
	"log"
	"strconv"
)

func NewRedisConn() (redis.Conn, error) {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		return nil, err
	}
	return c, nil
}

func AddUrlToDb(conn redis.Conn, urlData UrlData) error {
	_, err := conn.Do("RPUSH", urlData.ShortenUrl, urlData.OriginalUrl, urlData.CreationTime, urlData.ExpireTime, urlData.UserId) //just a simple Conn.do now, might add todo afterward
	if err != nil {
		return err
	}

	return nil
}

func CheckKeyExist(conn redis.Conn, key string) (int64, error) {
	result, err := conn.Do("EXISTS", key)
	if err != nil {
		return -1, err
	}
	return result.(int64), err
}

func GetOriginalUrlAndCheckIfExpire(conn redis.Conn, shortenUrl string) (string, bool, error) {
	log.Println("getting originate url and checking if expired")
	log.Println(shortenUrl)
	results, err := redis.Values(conn.Do("Lrange", shortenUrl, 0, 2))
	if err != nil {
		return "", false, err
	}

	var originalUrl, creationTime string
	var expiredTime int

	for i, v := range results {
		switch i {
		case 0:
			originalUrl = string(v.([]byte))
		case 1:
			creationTime = string(v.([]byte))
		case 2:
			expiredTime, err = strconv.Atoi(string(v.([]byte)))
			if err != nil {
				return "", true, err //consider wrong format also as expired, should be delete
			}
		default:
			return "", true, errors.New("incomplete url data") //consider incomplete data also as expired, should be delete
		}
	}

	elapsedTime := GetElapsedTime(creationTime)
	log.Println("this shorten url is posted for: ", elapsedTime)
	if elapsedTime > int64(expiredTime) {
		return "", true, nil
	}
	return originalUrl, false, nil
}

func DeleteUrlData(conn redis.Conn, shortenUrl string) error {
	_, err := conn.Do("DEL", shortenUrl) //not necessary to check if key is not exist
	if err != nil {
		return err
	}
	return nil
}
