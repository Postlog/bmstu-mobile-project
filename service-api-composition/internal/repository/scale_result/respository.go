package scale_result

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) Get(ctx context.Context, taskID string) (ScaleResult, error) {
	row := r.db.QueryRowContext(ctx, `select origin_image_id, scale_factor, image_id, error from scale_result where task_id = $1`, taskID)

	result := ScaleResult{TaskID: taskID}

	var imageID, errorText sql.NullString
	if err := row.Scan(&result.OriginalImageID, &result.ScaleFactor, &imageID, &errorText); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ScaleResult{}, ErrResultNotFound
		}

		return ScaleResult{}, fmt.Errorf("query row: %w", err)
	}

	if imageID.Valid {
		result.ImageID = &imageID.String
	} else if errorText.Valid {
		result.ErrorText = &errorText.String
	} else {
		return ScaleResult{}, fmt.Errorf("both image_id and error are null, taskID: %s", taskID)
	}

	return result, nil
}

func (r Repository) Save(ctx context.Context, results []ScaleResult) error {
	if len(results) == 0 {
		return nil
	}

	args := make([]interface{}, 0, len(results)*3)
	values := make([]string, 0, len(results))
	for i, req := range results {
		args = append(args, req.TaskID, req.OriginalImageID, req.ScaleFactor, req.ImageID, req.ErrorText)
		values = append(values,
			fmt.Sprintf(
				"($%d, $%d, $%d, $%d, $%d)",
				i*5+1, i*5+2, i*5+3, i*5+4, i*5+5,
			),
		)
	}

	query := `insert into scale_result (task_id, origin_image_id, scale_factor, image_id, error) values ` + strings.Join(values, ",")

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
