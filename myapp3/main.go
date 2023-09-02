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

func pathHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		// handle not found
		// w.WriteHeader(http.StatusNotFound)
		// fmt.Fprint(w, "Page not found")
		http.Error(w, "Page not found", http.StatusNotFound)
	}

	// if r.URL.Path == "/" {
	// 	homeHandler(w, r)
	// 	return
	// } else if r.URL.Path == "/contact" {
	// 	contactHandler(w, r)
	// 	return
	// }

	// handle not found

}

func main() {
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", handlerFunc)
	// http.HandleFunc("/", homeHandler)
	http.HandleFunc("/", pathHandler)

	// http.HandleFunc("/contact", contactHandler)
	/*
		경로 설정, DB 연결, 기타 수행 해야 하는 모든 작업에 대한 코드가 위치할 장소
	*/

	fmt.Println(":3000 포트에서 서버가 실행 중 입니다.")
	// fmt.Fprintln(os.Stdout, ":3000 !!")
	http.ListenAndServe(":3000", nil)
	// http.ListenAndServe(":3000", mux)
}
