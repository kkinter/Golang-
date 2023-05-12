package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"math/rand"
	"time"

	"greenlight.wook.net/internal/validator"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

// 구조체 태그를 추가하여 JSON으로 인코딩할 때 구조체가 표시되는 방식을 제어합니다.
type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	// 사용자 ID, 만료, 범위 정보가 포함된 토큰 인스턴스를 생성합니다. 만료 시간을
	//구하기 위해 현재 시간에 제공된 ttl(유효 기간) 기간 매개변수를 추가한 것이 보이시나요?
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	// 길이가 16바이트인 0 값의 바이트 슬라이스를 초기화합니다.
	randomBytes := make([]byte, 16)

	// crypto/rand 패키지의 Read() 함수를 사용하여 운영 체제의 CSPRNG에서
	// 임의의 바이트로 바이트 슬라이스를 채웁니다. CSPRNG가 제대로 작동하지 않으면 오류가 반환됩니다.
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// 바이트 슬라이스를 기본 32로 인코딩된 문자열로 인코딩하고 토큰 일반 텍스트 필드에 할당합니다.
	// 이 문자열은 환영 이메일에서 사용자에게 보내는 토큰 문자열이 됩니다. 다음과 비슷하게 보일 것입니다:
	//
	// y3qmgx3pj3wlrl2yrtqgq6krhu
	//
	// 기본적으로 기본 32 문자열은 끝에 = 문자가 추가될 수 있습니다.
	//토큰의 목적상 이 패딩 문자는 필요하지 않으므로 아래 줄에서
	// WithPadding(base32.NoPadding) 메서드를 사용하여 이를 생략합니다.

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	// 일반 텍스트 토큰 문자열의 SHA-256 해시를 생성합니다.
	// 이 값이 데이터베이스 테이블의 `hash` 필드에 저장됩니다. sha256.Sum256() 함수는
	// 길이 32의 *array*을 반환하므로, 작업하기 쉽도록 저장하기 전에 [:] 연산자를 사용하여
	// 슬라이스로 변환합니다.
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

// 일반 텍스트 토큰이 제공되었는지, 길이가 정확히 26바이트인지 확인합니다.
func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

// 토큰 모델 유형을 정의합니다.
type TokenModel struct {
	DB *sql.DB
}

// New() 메서드는 새 토큰 구조체를 생성한 다음 토큰 테이블에 데이터를 삽입하는 바로 가기입니다.
func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)

	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

// Insert()는 특정 토큰에 대한 데이터를 토큰 테이블에 추가합니다.
func (m TokenModel) Insert(token *Token) error {
	query := `
			INSERT INTO tokens (hash, user_id, expiry, scope)
			VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

// DeleteAllForUser()는 특정 사용자 및 범위에 대한 모든 토큰을 삭제합니다.
func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `
			DELETE FROM tokens
			WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err
}
