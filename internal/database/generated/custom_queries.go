package database

import "context"

const listCustomersPaginated = `-- name: ListCustomersPaginated :many
SELECT
    id,
    name,
    email,
    password,
    created_at,
    updated_at
FROM customers
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListCustomersPaginatedParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) ListCustomersPaginated(ctx context.Context, arg ListCustomersPaginatedParams) ([]Customer, error) {
	rows, err := q.db.Query(ctx, listCustomersPaginated, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Customer, 0)
	for rows.Next() {
		var i Customer
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Password,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
