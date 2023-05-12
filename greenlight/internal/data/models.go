package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit confilct")
)

// MovieModel을 감싸는 Models 구조체를 생성합니다.
// 빌드가 진행됨에 따라 UserModel 과 PermissionModel과 같은 다른 모델을 추가할 것입니다.
type Models struct {
	Movies interface {
		Insert(movie *Movie) error
		Get(id int64) (*Movie, error)
		Update(movie *Movie) error
		Delete(id int64) error
		GetAll(title string, genres []string, filters Filters) ([]*Movie, error)
	}
}

// 사용 편의성을 위해 초기화된 MovieModel을
// 포함하는 Models 구조체를 반환하는 New() 메서드도 추가했습니다.
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		Movies: MockMovieModel{},
	}
}
