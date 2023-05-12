package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"greenlight.wook.net/internal/validator"
)

// 사용자 정의 ErrDuplicateEmail 오류를 정의합니다.
var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

// 개별 사용자를 나타내는 User 구조체를 정의합니다.
// 중요한 것은 json:"-" 구조체 태그를 사용하여 JSON으로 인코딩할
// 때 비밀번호 및 버전 필드가 출력에 나타나지 않도록 하는 것입니다.
// 또한 Password 필드에 아래에 정의된 사용자 지정 비밀번호 유형을 사용하는 것도 주목하세요.
type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

// 연결 풀을 감싸는 UserModel 구조체를 생성합니다.
type UserModel struct {
	DB *sql.DB
}

// 사용자에 대한 비밀번호의 일반 텍스트 및 해시 버전을 포함하는
// 구조체인 사용자 정의 비밀번호 유형을 만듭니다. 일반 텍스트 필드는
// 문자열에 대한 *포인터*이므로 구조체에 전혀 존재하지 않는
// 일반 텍스트 비밀번호와 빈 문자열 ""인 일반 텍스트 비밀번호를 구분할 수 있습니다.
type password struct {
	plaintext *string
	hash      []byte
}

// 데이터베이스에 사용자에 대한 새 레코드를 삽입합니다.
// ID, created_at 및 버전 필드는 모두 데이터베이스에서 자동으로 생성되므로
// movie을 생성할 때와 같은 방식으로 삽입 후 RETURNING 절을 사용하여 사용자 구조체로 읽어들입니다.
func (m UserModel) Insert(user *User) error {

	query := `
		INSERT INTO users (name, email, password_hash, activated)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []any{user.Name, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 테이블에 이미 이 이메일 주소가 포함된 레코드가 있는 경우
	// 삽입을 수행하려고 하면 이전 장에서 설정한 UNIQUE "users_email_key"
	// 제약 조건을 위반하게 됩니다. 이 오류를 구체적으로 확인하고
	// 사용자 지정 ErrDuplicateEmail 오류를 반환합니다.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

// 사용자의 이메일 주소를 기준으로 데이터베이스에서 사용자 세부 정보를 검색합니다.
// 이메일 열에 UNIQUE 제약 조건이 있으므로 이 SQL 쿼리는 하나의 레코드만 반환합니다
// (또는 전혀 반환하지 않으며, 이 경우 ErrRecordNotFound 오류를 반환합니다).
func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, created_at, name, email, password_hash, activated, version
		FROM users
		WHERE email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

// 특정 사용자에 대한 세부 정보를 업데이트합니다. 동영상을 업데이트할 때와 마찬가지로
// 요청 주기 동안 경합 조건을 방지하기 위해 버전 필드에 대해 확인합니다. 또한
// 업데이트를 수행할 때 원래 사용자 레코드를 삽입할 때와 마찬가지로
// "users_email_key" 제약 조건 위반 여부도 확인합니다.
func (m UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version`

	args := []any{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail

		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Set() 메서드는 일반 텍스트 비밀번호의 bcrypt 해시를 계산하고
// 해시와 일반 텍스트 버전을 모두 구조체에 저장합니다.
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil

}

// Matches() 메서드는 제공된 일반 텍스트 암호가 구조체에 저장된 해시된 암호와
// 일치하는지 확인하여 일치하면 true를 반환하고 그렇지 않으면 false를 반환합니다.
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	// 독립 실행형 ValidateEmail() 도우미를 호출합니다.
	ValidateEmail(v, user.Email)

	// 일반 텍스트 비밀번호가 nil이 아닌 경우 독립 실행형
	// ValidatePasswordPlaintext() 헬퍼를 호출합니다.
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	// 비밀번호 해시가 0이면 코드베이스의 논리 오류 때문일 수 있습니다
	// (사용자 비밀번호를 설정하는 것을 잊어버렸기 때문일 수 있습니다).
	// 여기에 포함시키는 것은 유용한 유효성 검사이지만 클라이언트가 제공한 데이터에는 문제가 없습니다.
	// 따라서 유효성 검사 맵에 오류를 추가하는 대신 패닉을 발생시킵니다.
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
