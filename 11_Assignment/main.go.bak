package main

import (
	"fmt"
	"net/http"
)


func bar( res http.ResponseWriter, req *http.Request){
  res.Header().Set("Content-Type", "text/plain")
  fmt.Fprintf(res, "req.URL: %v\n", req.URL)
}

func main(){
  http.HandleFunc("/name/", bar)
  http.ListenAndServe(":8080", nil)
}
