package main

import (
	"flag"
	"log"
	"net/http"
)

// 웹페이지를 위한 HTML을 포함하는 문자열 상수를 정의합니다.
// 이는 <h1> 헤더 태그와, POST /v1/tokens/authentication 엔드포인트를 호출하고
//
//	응답 본문을 <div id="output"></div> 태그 안에 작성하는 JavaScript로 구성됩니다.
const html = `
<!DOCTYPE html>
<html lang="en">
   <head>
      <meta charset="UTF-8">
   </head>
   <body>
      <h1>Preflight CORS</h1>
      <div id="output"></div>
      <script>
         document.addEventListener('DOMContentLoaded', function() {
         fetch("http://localhost:4000/v1/tokens/authentication", {
         method: "POST",
         headers: {
         'Content-Type': 'application/json'
         },
         body: JSON.stringify({
         email: 'test@example.com',
         password: 'password'
         })
         }).then(
         function (response) {
         response.text().then(function (text) {
         document.getElementById("output").innerHTML = text;
         });
         },
         function(err) {
         document.getElementById("output").innerHTML = err;
         }
         );
         });
      </script>
   </body>
</html>`

func main() {
	addr := flag.String("addr", ":9000", "Server address")
	flag.Parse()
	log.Printf("starting server on %s", *addr)
	err := http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}))
	log.Fatal(err)
}
