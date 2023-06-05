package comment

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrFetchingComment = errors.New("failed to fectch comment by id")
	ErrNotImplemented  = errors.New("not implemented")
)

type Comment struct {
	ID     string
	Slug   string
	Body   string
	Author string
}

// Store - 이 인터페이스는 서비스가 작동하기 위해 필요한 모든 메소드를 정의합니다.
type Store interface {
	GetComment(context.Context, string) (Comment, error)
}

// 서비스 - 우리의 모든 로직이 구축되는 구조.
type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetComment(ctx context.Context, id string) (Comment, error) {
	fmt.Println("retrieving a comment")

	cmt, err := s.Store.GetComment(ctx, id)
	if err != nil {
		fmt.Println(err)
		return Comment{}, ErrFetchingComment
	}

	return cmt, nil
}

func (s *Service) UpdateComment(ctx context.Context, cmt Comment) error {
	return ErrNotImplemented
}

func (s *Service) DeleteComment(ctx context.Context, id string) error {
	return ErrNotImplemented
}

func (s *Service) CreateComment(ctx context.Context, cmt Comment) (Comment, error) {
	return Comment{}, ErrNotImplemented
}
