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

func convertCommentRowToComment(c CommentRow) comment.Comment {
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

	return convertCommentRowToComment(cmtRow), nil
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

func (d *Database) DeleteComment(ctx context.Context, id string) error {
	_, err := d.Client.ExecContext(
		ctx,
		`DELETE FROM comments where id = $1`,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to delete comment from db: %w", err)
	}

	return nil
}

func (d *Database) UpdateComment(ctx context.Context, id string, cmt comment.Comment) (comment.Comment, error) {
	query := `
		UPDATE comments SET
		slug = $1,
		author = $2,
		body = $3 
		WHERE id = $4
	`

	_, err := d.Client.ExecContext(ctx, query, cmt.Slug, cmt.Author, cmt.Body, id)
	if err != nil {
		return comment.Comment{}, fmt.Errorf("failed to update comment: %w", err)
	}

	return convertCommentRowToComment(CommentRow{
		ID:     id,
		Slug:   sql.NullString{String: cmt.Slug, Valid: true},
		Body:   sql.NullString{String: cmt.Body, Valid: true},
		Author: sql.NullString{String: cmt.Author, Valid: true},
	}), nil
}
