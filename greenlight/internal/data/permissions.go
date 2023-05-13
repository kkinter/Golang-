package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Permissions []string

// 단일 사용자에 대한 권한 코드(예: "movies:read" 및 "movies:write")를
// 보관하는 데 사용할 권한 슬라이스를 정의합니다.
func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

type PermissionModel struct {
	DB *sql.DB
}

// GetAllForUser() 메서드는 권한 슬라이스에서 특정 사용자에 대한 모든 권한 코드를 반환합니다.
// 이 메서드의 코드는 매우 친숙하게 느껴질 것입니다. SQL 쿼리에서 여러 데이터 행을 검색하는
// 데 이미 보았던 표준 패턴을 사용합니다.
func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	query := `
		SELECT permissions.code
		FROM permissions
		INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
		INNER JOIN users ON users_permissions.user_id = users.id
		WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions
	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return permissions, nil
}

// 특정 사용자에 대해 제공된 권한 코드를 추가합니다. 한 번의 호출로 여러
// 권한을 할당할 수 있도록 코드에 가변 매개 변수를 사용하고 있다는 점에 유의하세요.

func (m PermissionModel) AddForUser(userID int64, codes ...string) error {
	query := `
			INSERT INTO users_permissions
			SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userID, pq.Array(codes))
	return err
}
