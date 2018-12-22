package todo

import (
	"context"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"main/models"
)

type App struct {
	db boil.ContextExecutor
}

type TodoAndUser struct {
	models.Todo `boil:",bind"`
	models.User `boil:",bind"`
}

func (a *App) FetchTaskAndUser() (*models.User, error) {
	task, _ := models.FindTodo(context.Background(), a.db, 1)
	user, err := task.User().One(context.Background(), a.db)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *App) FetchTasksAndUser() ([]*TodoAndUser, error) {
	var tau []*TodoAndUser
	err := models.NewQuery(
		qm.Select("todo.id", "todo.title", "todo.note", "user.name"),
		qm.From("todo"),
		qm.InnerJoin("user on todo.user_id = user.id"),
	).Bind(context.Background(), a.db, &tau)

	return tau, err
}

func (a *App) FetchUnfinished() ([]*models.Todo, error) {
	queries := []qm.QueryMod{
		qm.Where(models.TodoColumns.Finished+"=?", false),
	}

	counts, err := models.Todos(queries...).Count(context.Background(), a.db)
	if err != nil {
		return nil, err
	} else if counts == 0 {
		return make([]*models.Todo, 0), nil
	}

	// fetch!
	return models.Todos(queries...).All(context.Background(), a.db)
}

func (a *App) Store(obj []*models.Todo) error {
	// insert for each todo
	for _, v := range obj {
		if err := v.Insert(context.Background(), a.db, boil.Infer()); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Finish(ids []int64) error {
	// Set updating columns
	updCols := map[string]interface{}{
		models.TodoColumns.Finished: true,
	}

	// WhereIn method needs to pass a slice of interface{}
	targetIDs := make([]interface{}, len(ids))
	for i, v := range ids {
		targetIDs[i] = v
	}
	query := qm.WhereIn(models.TodoColumns.ID+" IN ?", targetIDs...)

	// update!
	_, err := models.Todos(query).UpdateAll(context.Background(), a.db, updCols)

	return err
}
