package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"greenlight.wook.net/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

type MovieModel struct {
	DB *sql.DB
}

// Insert() 메서드는 새 레코드에 대한
// 데이터를 포함해야 하는 movie 구조체에 대한 포인터를 받습니다.
func (m MovieModel) Insert(movie *Movie) error {

	// 영화 테이블에 새 레코드를 삽입하고
	// 시스템에서 생성된 데이터를 반환하기 위한 SQL 쿼리를 정의합니다.
	query := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`
	// movie 구조체에서 플레이스홀더 매개변수의 값을 포함하는 args 슬라이스를 생성합니다.
	// 이 슬라이스를 SQL 쿼리 바로 옆에 선언하면 쿼리에서
	// *어떤 값이 어디에 사용되는지* 명확하게 파악할 수 있습니다.
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	// QueryRow() 메서드를 사용하여 연결 풀에서 SQL 쿼리를 실행하고,
	// 가변 파라미터로 args 슬라이스를 전달하고,
	// 시스템에서 생성된 id, created_at 및 버전 값을 movie 구조체로 스캔합니다.
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {

	// movie  ID에 사용하는 PostgreSQL bigserial  유형은 기본적으로 1에서 자동 증가를 시작하므로
	//  이보다 작은 ID 값을 갖는 영화는 없다는 것을 알 수 있습니다.
	// 불필요한 데이터베이스 호출을 피하기 위해 바로 가기를 사용하여
	// ErrRecordNotFound 오류를 바로 반환합니다.
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE id = $1`

	// 쿼리에서 반환된 데이터를 저장할 Movie 구조체를 선언합니다.
	var movie Movie

	// QueryRow() 메서드를 사용하여 쿼리를 실행하고,
	// 제공된 id 값을 자리 표시자 매개변수로 전달한 다음,
	// 응답 데이터를 Movie 구조체의 필드로 스캔합니다.
	//  중요한 점은 pq.Array() 어댑터 함수를 사용하여 장르 열의 스캔 대상을
	// 다시 변환해야 한다는 점입니다.
	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	// 오류를 처리합니다. 일치하는 동영상을 찾지 못하면 Scan()은 sql.ErrNoRows 오류를
	// 반환합니다. 이 오류를 확인하고 대신 커스텀 ErrRecordNotFound 오류를 반환합니다.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	query := `
		UPDATE movies
		SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
		WHERE id = $5
		RETURNING version `

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	// 쿼리를 실행하기 위해 QueryRow() 메서드를 사용하고,
	// 가변 파라미터로 args 슬라이스를 전달하고, 새 버전 값을 movie  구조체로 스캔합니다.
	return m.DB.QueryRow(query, args...).Scan(&movie.Version)
}

func (m MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM movies
		WHERE id = $1`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

type MockMovieModel struct{}

func (m MockMovieModel) Insert(movie *Movie) error {
	return nil
}
func (m MockMovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}
func (m MockMovieModel) Update(movie *Movie) error {
	return nil
}
func (m MockMovieModel) Delete(id int64) error {
	return nil
}
