package main

import (
	"github.com/russross/blackfriday"
	"net/http"
	"fmt"
	"strconv"
	"time"
	"strings"
)

func ConvertMarkdownToHtml(markdawn string) string {
	return string(blackfriday.MarkdownBasic([]byte(markdawn)))
}

func CheckCookies(r *http.Request) bool {
	a := false
	for _, cookie := range r.Cookies() {
	    if cookie.Value == config.Pass {
	    	a = true
	    } else {
	    	a = false
	    }
	}
	return a
}

func timeX(t string) string{
	i, _ := strconv.ParseInt(t, 10, 64)
	tm := time.Unix(i, 0)
	Time := tm.Format("_2-01-2006")
	return Time
}

func log(STR string) {
	fmt.Printf("\x1b[32;1m%s\x1b[0m\n", STR)
}

func tagsX(t string, u string) string{
	t_ := strings.Split(t, ", ")
	var tags string
	var ta string
	for _,element := range t_ {
		if element == "" {
			continue
		}
		if u != "true" {
	      	ta += element+", "
	      	tagsz := len(ta)
	      	tags = ta[:tagsz-2]
	   	} else {
	   		tags += "<a class='tag' href='/tag/"+element+"'>"+element+"</a>"
	   	}
	}
	return tags
}

func timeRFC(t string) string{
	i, _ := strconv.ParseInt(t, 10, 64)
	tm := time.Unix(i, 0)
	Time := tm.Format(time.RFC1123Z)
	return Time
}