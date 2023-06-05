package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-api-prac/internal/comment"
	uuid "github.com/satori/go.uuid"
)

type CommentRow struct {
	ID     string
	Slug   sql.NullString
	Body   sql.NullString
	Author sql.NullString
}

func converCommentRowToComment(c CommentRow) comment.Comment {
	return comment.Comment{
		ID:     c.ID,
		Slug:   c.Slug.String,
		Author: c.Author.String,
		Body:   c.Body.String,
	}
}

func (d *Database) GetComment(ctx context.Context, uuid string) (comment.Comment, error) {
	var cmtRow CommentRow

	row := d.Client.QueryRowContext(
		ctx,
		`SELECT id, slug, body, author
		FROM comments
		WHERE id = $1`,
		uuid,
	)

	err := row.Scan(&cmtRow.ID, &cmtRow.Slug, &cmtRow.Body, &cmtRow.Author)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("error fetching the comment by uuid: %w", err)
	}

	return converCommentRowToComment(cmtRow), nil
}

func (d *Database) PostComment(ctx context.Context, cmt comment.Comment) (comment.Comment, error) {
	cmt.ID = uuid.NewV4().String()

	stmt, err := d.Client.PrepareContext(ctx, `
		INSERT INTO comments (id, slug, author, body)
		VALUES ($1, $2, $3, $4)
	`)

	if err != nil {
		return comment.Comment{}, fmt.Errorf("failed to prepare SQL statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, cmt.ID, cmt.Slug, cmt.Author, cmt.Body)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("failed to insert comment: %w", err)
	}

	return cmt, nil
}
