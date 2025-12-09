package controllers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type TaskController struct {
	taskService app.TaskService
}

func NewTaskController(ts app.TaskService) TaskController {
	return TaskController{
		taskService: ts,
	}
}

func (c TaskController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController.Save(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		task.UserId = user.Id
		task.Status = domain.NewTaskStatus
		task, err = c.taskService.Save(task)
		if err != nil {
			log.Printf("TaskController.Save(c.taskService.Save): %s", err)
			InternalServerError(w, err)
			return
		}

		taskDto := resources.TaskDto{}
		taskDto = taskDto.DomainToDto(task)
		Success(w, taskDto)
	}
}

func (c TaskController) FindList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		status := ""
		if r.URL.Query().Has("status") {
			status = r.URL.Query().Get("status")
		}

		search := ""
		if r.URL.Query().Has("search") {
			search = r.URL.Query().Get("search")
		}

		var deadlineFrom *time.Time
		var deadlineTo *time.Time
		if r.URL.Query().Has("date") {
			dateStr := r.URL.Query().Get("date")
			day, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				log.Printf("TaskController.FindList: invalid date param: %s", err)
				BadRequest(w, errors.New("invalid date format, expected YYYY-MM-DD"))
				return
			}
			start := day
			end := day.Add(24 * time.Hour)
			deadlineFrom = &start
			deadlineTo = &end
		}

		filters := database.TaskFilters{
			UserId:       user.Id,
			Status:       status,
			Search:       search,
			DeadlineFrom: deadlineFrom,
			DeadlineTo:   deadlineTo,
		}

		tasks, err := c.taskService.FindList(filters)
		if err != nil {
			log.Printf("TaskController.FindList(c.taskService.FindList): %s", err)
			InternalServerError(w, err)
			return
		}

		tasksDto := resources.TasksDto{}
		tasksDto = tasksDto.DomainToDto(tasks)
		Success(w, tasksDto)
	}
}

func (c TaskController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task := r.Context().Value(TaskKey).(domain.Task)

		if user.Id != task.UserId {
			err := errors.New("access denied")
			Forbidden(w, err)
			return
		}

		taskDto := resources.TaskDto{}
		taskDto = taskDto.DomainToDto(task)
		Success(w, taskDto)
	}
}

func (c TaskController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task := r.Context().Value(TaskKey).(domain.Task)

		if user.Id != task.UserId {
			err := errors.New("access denied")
			Forbidden(w, err)
			return
		}

		update, err := requests.Bind(r, requests.UpdateTaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController.Update(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		task.Status = update.Status
		task, err = c.taskService.Update(task)
		if err != nil {
			log.Printf("TaskController.Update(c.taskService.Update): %s", err)
			InternalServerError(w, err)
			return
		}

		taskDto := resources.TaskDto{}
		taskDto = taskDto.DomainToDto(task)
		Success(w, taskDto)
	}
}

func (c TaskController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task := r.Context().Value(TaskKey).(domain.Task)

		if user.Id != task.UserId {
			err := errors.New("access denied")
			Forbidden(w, err)
			return
		}

		if err := c.taskService.Delete(task.Id); err != nil {
			log.Printf("TaskController.Delete(c.taskService.Delete): %s", err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}
