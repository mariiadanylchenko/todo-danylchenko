package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const TasksTableName = "tasks"

type task struct {
	Id          uint64            `db:"id,omitempty"`
	UserId      uint64            `db:"user_id"`
	Title       string            `db:"title"`
	Description *string           `db:"description"`
	Status      domain.TaskStatus `db:"status"`
	Deadline    *time.Time        `db:"deadline"`
	CreatedDate time.Time         `db:"created_date"`
	UpdatedDate time.Time         `db:"updated_date"`
	DeletedDate *time.Time        `db:"deleted_date"`
}

type TaskRepository interface {
	Save(t domain.Task) (domain.Task, error)
	FindList(tf TaskFilters) ([]domain.Task, error)
	Find(id uint64) (domain.Task, error)
	Update(t domain.Task) (domain.Task, error)
	Delete(id uint64) error
}

type taskRepository struct {
	sess db.Session
	coll db.Collection
}

func NewTaskRepository(dbSession db.Session) TaskRepository {
	return taskRepository{
		sess: dbSession,
		coll: dbSession.Collection(TasksTableName),
	}
}

func (r taskRepository) Save(t domain.Task) (domain.Task, error) {
	tsk := r.mapDomainToModel(t)
	tsk.CreatedDate = time.Now()
	tsk.UpdatedDate = time.Now()

	err := r.coll.InsertReturning(&tsk)
	if err != nil {
		return domain.Task{}, err
	}

	t = r.mapModelToDomain(tsk)
	return t, nil
}

func (r taskRepository) Find(id uint64) (domain.Task, error) {
	var t task

	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&t)
	if err != nil {
		return domain.Task{}, err
	}

	return r.mapModelToDomain(t), nil
}

func (r taskRepository) Update(t domain.Task) (domain.Task, error) {
	model := r.mapDomainToModel(t)
	model.UpdatedDate = time.Now()

	err := r.coll.Find(db.Cond{"id": model.Id, "deleted_date": nil}).Update(&model)
	if err != nil {
		return domain.Task{}, err
	}

	return r.mapModelToDomain(model), nil
}

func (r taskRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

type TaskFilters struct {
	UserId       uint64
	Status       string
	Search       string
	Deadline     *time.Time
	DeadlineFrom *time.Time
	DeadlineTo   *time.Time
}

func (r taskRepository) FindList(tf TaskFilters) ([]domain.Task, error) {
	var tasks []task
	query := r.coll.Find(db.Cond{"user_id": tf.UserId, "deleted_date": nil})

	if tf.Status != "" {
		query = query.And(db.Cond{"status": tf.Status})
	}

	if tf.Search != "" {
		query = query.And(db.Cond{"title ilike": "%" + tf.Search + "%"})
	}

	if tf.Deadline != nil {
		query = query.And(db.Cond{"deadline": tf.Deadline})
	}
	if tf.DeadlineFrom != nil {
		query = query.And(db.Cond{"deadline >=": *tf.DeadlineFrom})
	}
	if tf.DeadlineTo != nil {
		query = query.And(db.Cond{"deadline <": *tf.DeadlineTo})
	}

	err := query.All(&tasks)
	if err != nil {
		return nil, err
	}

	return r.mapModelToDomainCollection(tasks), nil
}

func (r taskRepository) mapDomainToModel(t domain.Task) task {
	return task{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Deadline:    t.Deadline,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}

func (r taskRepository) mapModelToDomain(t task) domain.Task {
	return domain.Task{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Deadline:    t.Deadline,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}

func (r taskRepository) mapModelToDomainCollection(ts []task) []domain.Task {
	tasks := make([]domain.Task, len(ts))
	for i, t := range ts {
		tasks[i] = r.mapModelToDomain(t)
	}
	return tasks
}
