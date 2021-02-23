package worker

import (
	"time"
)

type Job struct {
	name     string
	duration time.Duration
	jobFunc  func() error
}

type Jobs struct {
	j map[string]*Job
}

func (jobs *Jobs) executeJob(name string) error {

	err := jobs.j[name].jobFunc()
	if err != nil {
		return err
	}
	return nil
}

func createJob(name string, jobFunc func() error) {
	job := &Job{name, nil, jobFunc}
	jobs := &Jobs{make(map[string]*Job)}
	jobs.j[name] = job
}

func CreateWorkersForJob(name string, maxQueueSize, maxWorkers, port int) {
	var Jobs *Jobs

	for i := 1; i <= maxWorkers; i++ {
		go func(i int) {

			err := Jobs.executeJob(name)
			if err != nil {
				panic(err)
			}

		}(i)
	}
}
