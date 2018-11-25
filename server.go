package fetchtask

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskId string
type UniqueIdGenerator func() TaskId

type Request struct {
	Method  string
	Path    string
	Host    string
	Body    string
	Headers map[string]string
}

var TaskNotFound error

type FetchTask struct {
	Id      TaskId
	Request Request
	Error   error
}

type Server struct {
	generateId UniqueIdGenerator
	tasks      []FetchTask
}

func getTaskById(arr []FetchTask, id TaskId) (*FetchTask, int, error) {
	for i, task := range arr {
		if task.Id == id {
			return &task, i, nil
		}
	}
	return nil, 0, TaskNotFound
}

func newTaskFactory(generateId UniqueIdGenerator) func(r Request) FetchTask {
	return func(r Request) FetchTask {
		return FetchTask{
			Id:      generateId(),
			Request: r,
		}
	}
}

func executeRequest(r Request) error {
	return nil
}

func NewIdGeneratorMock() UniqueIdGenerator {
	i := 0
	return func() TaskId {
		i++
		return TaskId(fmt.Sprintf("%016d", i))
	}
}

func NewServer(id UniqueIdGenerator) *Server {
	return &Server{id, []FetchTask{}}
}

func (s *Server) Listen() {
	router := gin.Default()
	newTask := newTaskFactory(s.generateId)
	router.POST("/task", func(ctx *gin.Context) {
		var r Request
		if err := ctx.ShouldBindJSON(&r); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		task := newTask(r)
		task.Error = executeRequest(r)
		s.tasks = append(s.tasks, task)
		ctx.JSON(200, gin.H{
			"id": task.Id,
		})

	})
	router.GET("/tasks", func(ctx *gin.Context) {
		ctx.JSON(200, s.tasks)
	})
	router.GET("/task/:id", func(ctx *gin.Context) {
		if task, _, err := getTaskById(s.tasks, TaskId(ctx.Param("id"))); err != nil {
			ctx.JSON(404, gin.H{
				"error": "task not found",
			})
		} else {
			ctx.JSON(200, task)
		}
	})
	router.DELETE("/task/:id", func(ctx *gin.Context) {
		if _, index, err := getTaskById(s.tasks, TaskId(ctx.Param("id"))); err != nil {
			ctx.JSON(404, gin.H{
				"error": "task not found",
			})
		} else {
			s.tasks = append(s.tasks[:index], s.tasks[index+1:]...)
			ctx.Status(200)
		}

	})
	router.Run()
}
