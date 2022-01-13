package services

import (
	"errors"
	"fmt"
	"sync"

	"github.com/pete-robinson/go-cli-test/internal/domain/model"
)

const maxWorkers = 5

// repo interface
type Repository interface {
	FindById(int) (*model.Todo, error)
	Find(int) ([]*model.Todo, error)
}

// todo service
type TodoService struct {
	repository Repository
}

// create new todo service
func NewTodoService(repository Repository) *TodoService {
	return &TodoService{
		repository: repository,
	}
}

// find all
func (s *TodoService) FetchAll(numResults int) ([]*model.Todo, error) {
	res, err := s.repository.Find(numResults)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// fetch one
func (s *TodoService) FetchOne(id int) (*model.Todo, error) {
	res, err := s.repository.FindById(id)
	if err != nil {
		return nil, err
	}

	if res.Id == 0 {
		return nil, errors.New(fmt.Sprintf("No results found for id: %d", id))
	}

	return res, nil
}

// fetch list
func (s *TodoService) FetchList(list string) ([]*model.Todo, error) {
	// read and parse list
	resp, err := readFile(list)
	if err != nil {
		return nil, err
	}

	// fetch results
	results, err := s.fetchWithWorker(resp)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// fetch many individual rows by ID
// using a worker pool to speed up response time
func (s *TodoService) fetchWithWorker(ids *ReaderResponse) ([]*model.Todo, error) {
	numJobs := len(*ids)                         // number of results to fetch
	chJobs := make(chan int, numJobs)            // jobs channel
	chResults := make(chan *model.Todo, numJobs) // results channel
	chErrors := make(chan error, numJobs)        // error channel so we can catch errors across threads

	// init the worker pool
	var wg sync.WaitGroup
	for w := 0; w <= maxWorkers; w++ {
		wg.Add(1)
		go s.worker(chJobs, chResults, chErrors, &wg)
	}

	// defer closure of channels
	go func() {
		defer close(chResults)
		defer close(chErrors)
		wg.Wait()
	}()

	// pass IDs to jobs channel to init workers
	for _, i := range *ids {
		chJobs <- i
	}
	close(chJobs)

	// more robust error management needed
	// check for errors first and return if one is encountered
	var errors []error
	for e := range chErrors {
		errors = append(errors, e)
	}

	if len(errors) > 0 {
		return nil, errors[0]
	}

	// fetch results from chanel and return
	var res []*model.Todo
	for r := range chResults {
		res = append(res, r)
	}

	return res, nil
}

func (s *TodoService) worker(jobs <-chan int, results chan<- *model.Todo, errors chan<- error, wg *sync.WaitGroup) {
	// close the wait group for this worker once worker is complete
	defer wg.Done()

	// loop over jobs for this worker
	for j := range jobs {
		// fetch result from API
		res, err := s.FetchOne(j)
		if err != nil {
			// send error to errors channel
			errors <- err
		} else {
			// send results to results channel
			results <- res
		}
	}
}
