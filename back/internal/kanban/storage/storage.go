package storage

import (
	"context"
	"database/sql"
	"github.com/aakosarev/kanban-board/back/internal/kanban"
	"github.com/aakosarev/kanban-board/back/internal/models"
	"github.com/aakosarev/kanban-board/back/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type KanbanStorage struct {
	log    logger.Logger
	client *pgxpool.Pool
}

func NewKanbanStorage(log logger.Logger, client *pgxpool.Pool) kanban.Storage {
	return &KanbanStorage{
		log:    log,
		client: client,
	}
}

func (k *KanbanStorage) CreateColumn(ctx context.Context, column *models.Column) (*models.Column, error) {
	query := `
		INSERT INTO "column"(user_id, name)
		VALUES ($1, $2)
		RETURNING *;
	`

	c := &models.Column{}

	if err := k.client.QueryRow(ctx, query, column.UserID, column.Name).Scan(&c.ID, &c.UserID, &c.Name); err != nil {
		return nil, err
	}

	return c, nil
}

func (k *KanbanStorage) DeleteColumn(ctx context.Context, id int) error {
	query := `
		DELETE FROM "column"
		WHERE id = $1;
	`

	res, err := k.client.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "KanbanStorage.DeleteColumn.Exec")
	}

	if rowsAffected := res.RowsAffected(); rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "KanbanStorage.DeleteColumn.rowsAffected")
	}

	return nil
}

func (k *KanbanStorage) ChangeNameColumn(ctx context.Context, column *models.Column) (*models.Column, error) {
	query := `
		UPDATE "column" 
		SET name = $1
		WHERE id = $2
		RETURNING *;
	`

	c := &models.Column{}

	if err := k.client.QueryRow(ctx, query, column.Name, column.ID).Scan(&c.ID, &c.UserID, &c.Name); err != nil {
		return nil, err
	}

	return c, nil
}

func (k *KanbanStorage) CreateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	query := `
		INSERT INTO "task"(column_id, description)
		VALUES ($1, $2)
		RETURNING *;
	`

	t := &models.Task{}

	if err := k.client.QueryRow(ctx, query, task.ColumnID, task.Description).Scan(&t.ID, &t.ColumnID, &t.Description); err != nil {
		return nil, err
	}

	return t, nil
}

func (k *KanbanStorage) DeleteTask(ctx context.Context, id int) error {
	query := `
		DELETE FROM "task"
		WHERE id = $1;
	`

	res, err := k.client.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "KanbanStorage.DeleteTask.Exec")
	}

	if rowsAffected := res.RowsAffected(); rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "KanbanStorage.DeleteTask.rowsAffected")
	}

	return nil
}

func (k *KanbanStorage) ChangeDescriptionTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	query := `
		UPDATE "task" 
		SET description = $1
		WHERE id = $2
		RETURNING *;
	`

	t := &models.Task{}

	if err := k.client.QueryRow(ctx, query, task.Description, task.ID).Scan(&t.ID, &t.ColumnID, &t.Description); err != nil {
		return nil, err
	}

	return t, nil
}

func (k *KanbanStorage) ChangeColumnIDTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	query := `
		UPDATE "task" 
		SET column_id = $1
		WHERE id = $2
		RETURNING *;
	`

	t := &models.Task{}

	if err := k.client.QueryRow(ctx, query, task.ColumnID, task.ID).Scan(&t.ID, &t.ColumnID, &t.Description); err != nil {
		return nil, err
	}

	return t, nil
}

func (k *KanbanStorage) GetKanbanBoardByUserID(ctx context.Context, userID int) (*models.Board, error) {
	query := `
		SELECT 
		    "column".id AS column_id, 
		    "column".name AS column_name,
		    "task".id AS task_id,
		    "task".column_id AS task_column_id,
		    "task".description AS task_description
		FROM "column"
		LEFT JOIN "task" ON "column".id = "task".column_id
		WHERE "column".user_id = $1
		ORDER BY column_id;
	`

	rows, err := k.client.Query(ctx, query, userID)
	if err != nil {
		return nil, errors.Wrap(err, "KanbanStorage.GetKanbanBoardByUserID.Query")
	}
	defer rows.Close()

	b := models.Board{}
	columnsMap := make(map[int32]*models.Col)

	for rows.Next() {
		k.log.Debug(1)
		var colID, taskID, taskColumnID sql.NullInt32
		var colName, taskDesc sql.NullString
		if err := rows.Scan(&colID, &colName, &taskID, &taskColumnID, &taskDesc); err != nil {
			return nil, errors.Wrap(err, "KanbanStorage.GetKanbanBoardByUserID.Scan")
		}
		col, exists := columnsMap[colID.Int32]
		if !exists {
			col = &models.Col{
				ID:   int(colID.Int32),
				Name: colName.String,
			}
			columnsMap[colID.Int32] = col
			b.Columns = append(b.Columns, col)
		}

		task := models.T{
			ID:          int(taskID.Int32),
			ColumnID:    int(taskColumnID.Int32),
			Description: taskDesc.String,
		}
		col.Tasks = append(col.Tasks, &task)

	}
	k.log.Debug(b)
	return &b, nil
}
