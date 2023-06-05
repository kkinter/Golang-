package http

// func JWTAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := r.Header["Authorization"]
// 		if authHeader == nil {
// 			http.Error(w, "not auth", http.StatusUnauthorized)
// 			return
// 		}

// 		authHeaderParts := strings.Split(authHeader[0], " ")
// 		if len(authHeaderParts) != 2 || strings.ToLower(authHeader[0]) != "bearer" {
// 			http.Error(w, "not auth", http.StatusUnauthorized)
// 			return
// 		}

// 		if validateToken(authHeaderParts[1]) {
// 			original(w, r)
// 		} else {
// 			http.Error(w, "not auth", http.StatusUnauthorized)
// 			return
// 		}
// 	}
// }

// func validateToken(accessToken string) bool {
// 	var mySigningKey = []byte("wookseong")
// 	token, err := jwt.Par

// }
