package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "<h1>remind http pacsksage</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Contact Page</h1><p><small>ifol1129@gmail.com</small>")
}

// func pathHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/":
// 		homeHandler(w, r)
// 	case "/contact":
// 		contactHandler(w, r)
// 	default:
// 		// handle not found
// 		http.Error(w, "Page not found", http.StatusNotFound)
// 	}
// }

type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		// handle not found
		http.Error(w, "Page not found", http.StatusNotFound)
	}
}

func main() {
	var router Router
	fmt.Println(":3000 포트에서 서버가 실행 중 입니다.")
	http.ListenAndServe(":3000", router)

	// 1.
	// var router http.HandlerFunc = pathHandler
	// fmt.Println(":3000 포트에서 서버가 실행 중 입니다.")
	// http.ListenAndServe(":3000", router)

	// 2.
	// fmt.Println(":3000 포트에서 서버가 실행 중 입니다.")
	// http.ListenAndServe(":3000", http.HandlerFunc(pathHandler))

	// var a int64 = 123
	// var b int32
	// b = int32(a)
	// fmt.Println(a)
	// fmt.Println(b)

	// http.Handler = ServeHTTP 메서드가 있는 인터페이스
	// http.HandlerFunc = ServeHTTP 메서드와 동일한 인수를 받는 함수.
	// 또한 http.Handler를 구현함.
}
