package tasks

import (
	"github.com/monstorak/monstorak/pkg/client"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("tasks")

// TaskRunner holds a reference to a Client and a list of TaskSpec
type TaskRunner struct {
	client *client.Client
	tasks  []*TaskSpec
}

// NewTaskRunner creates a new TaskRunner object
func NewTaskRunner(client *client.Client, tasks []*TaskSpec) *TaskRunner {
	return &TaskRunner{
		client: client,
		tasks:  tasks,
	}
}

// RunAll executes a list of Tasks held by TaskRunner
func (tl *TaskRunner) RunAll() error {
	var g errgroup.Group

	for i, ts := range tl.tasks {
		// shadow vars due to concurrency
		ts := ts
		i := i

		g.Go(func() error {
			log.V(4).Info("running task %d of %d: %v", i+1, len(tl.tasks), ts.Name)
			err := tl.ExecuteTask(ts)
			log.V(4).Info(("ran task %d of %d: %v", i+1, len(tl.tasks), ts.Name)
			return errors.Wrapf(err, "running task %v failed", ts.Name)
		})
	}

	return g.Wait()
}

// ExecuteTask runs each individual task
func (tl *TaskRunner) ExecuteTask(ts *TaskSpec) error {
	return ts.Task.Run()
}

// NewTaskSpec returns a new TaskSpec object
func NewTaskSpec(name string, task Task) *TaskSpec {
	return &TaskSpec{
		Name: name,
		Task: task,
	}
}

//TaskSpec holds information of a particular Task
type TaskSpec struct {
	Name string
	Task Task
}

// Task interface which will be implemented by each sub-task
type Task interface {
	Run() error
}
