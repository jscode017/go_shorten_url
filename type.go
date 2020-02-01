//from https://www.educative.io/courses/grokking-the-system-design-interview/m2ygV4E81AR
package main

type UrlData struct {
	ShortenUrl   string
	OriginalUrl  string
	CreationTime string
	ExpireTime   string
	UserId       string // converted from int

}

type UserIdData struct {
}

//func (t Time) String() string
//func (t Time) Unix() int64
//2001-09-09 01:46:40 +0000 UTC
//func Parse("2001-09-09 01:46:40 +0000 UTC", value string) (Time, error)
