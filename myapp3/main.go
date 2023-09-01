package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>remind http package</h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	fmt.Println(":3000 포트에서 서버가 실행 중 입니다.")
	http.ListenAndServe(":3000", nil)
}
