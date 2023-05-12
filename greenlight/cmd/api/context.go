package main

import (
	"context"
	"net/http"

	"greenlight.wook.net/internal/data"
)

// 커스텀 contextKey 타입을 정의합니다. 기반 타입은 string입니다.
type contextKey string

// 문자열 "user"를 contextKey 타입으로 변환하고 userContextKey 상수에 할당합니다.
// 이 상수를 요청 컨텍스트에서 사용자 정보를 가져오고 설정하는 데 키로 사용할 것입니다.
const userContextKey = contextKey("user")

// contextSetUser() 메서드는 제공된 User 구조체를 컨텍스트에 추가한 새 요청 사본을 반환합니다.
// 여기서 우리는 userContextKey 상수를 키로 사용합니다.
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser()는 요청 컨텍스트에서 User 구조체를 가져옵니다.
// 이 헬퍼를 사용하는 경우에는 컨텍스트에 User 구조체 값이 있는 것을 논리적으로 예상하며,
// 그런 값이 없다면 '예기치 않은' 오류가 될 것입니다.
// 이 책에서 이전에 설명했듯이, 이러한 경우에는 패닉해도 괜찮습니다
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
