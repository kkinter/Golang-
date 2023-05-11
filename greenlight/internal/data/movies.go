package data

import (
	"time"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	// int32 대신 Runtime 타입을 사용합니다. Runtime 필드의 기본값이 0이면
	// 비어 있는 것으로 간주되어 생략되며, 방금 만든 MarshalJSON() 메서드는 호출되지 않습니다.
	Runtime Runtime  `json:"runtime,omitempty"`
	Genres  []string `json:"genres,omitempty"`
	Version int32    `json:"version"`
}
