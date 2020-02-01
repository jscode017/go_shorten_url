package main

import (
	"log"
	"testing"
	"time"
)

func TestDbOperations(t *testing.T) {
	conn, err := NewRedisConn()
	if err != nil {
		t.Fatal(err)
	}

	urlData := UrlData{
		OriginalUrl:  "www.google.com",
		CreationTime: ConvertTimeToStr(time.Now()),
		ExpireTime:   "5000",
		UserId:       "-1",
	}
	urlData.ShortenUrl = GenerateKey(urlData.OriginalUrl)

	err = AddUrlToDb(conn, urlData)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(urlData.ShortenUrl)

	keyExist, err := CheckKeyExist(conn, urlData.ShortenUrl)
	if err != nil {
		t.Fatal(err)
	}
	if keyExist != 1 {
		t.Fatal("add urlData failed")
	}

	originalUrl, ifExpired, err := GetOriginalUrlAndCheckIfExpire(conn, urlData.ShortenUrl)
	if err != nil {
		t.Fatal(err)
	}
	if ifExpired == true {
		t.Fatal("should not expired")
	}
	log.Println(originalUrl)

	err = DeleteUrlData(conn, urlData.ShortenUrl)
	if err != nil {
		t.Fatal("delete url Data error")
	}

	keyExist, err = CheckKeyExist(conn, urlData.ShortenUrl)
	if err != nil {
		t.Fatal(err)
	}
	if keyExist != 0 {
		t.Fatal("key should not exist")
	}
}
