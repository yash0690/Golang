
 package main

import (
  "fmt"
 	"io"
 	"net/http"
 	"strconv"
)

func thisLittleWebpage(res http.ResponseWriter, req *http.Request) {
 
  if req.URL.Path != "/" {
    http.NotFound(res, req)
    return
  }

  cookie, err := req.Cookie("my-cookie")
 
  if err == http.ErrNoCookie {
    cookie = &http.Cookie{
      Name:  "my-cookie",
      Value: "0",
    }
  }
 
  count, _ := strconv.Atoi(cookie.Value)
  count++ 
  cookie.Value = strconv.Itoa(count) 

  http.SetCookie(res, cookie)

  io.WriteString(res, cookie.Value)
}

func main() {
  http.HandleFunc("/", thisLittleWebpage)

  fmt.Println("server is now running...") 
  http.ListenAndServe(":8080", nil) 
}
