package repository

import (
	"encoding/json"
	"strconv"

	"github.com/pete-robinson/go-cli-test/internal/domain/model"
	"github.com/pete-robinson/go-cli-test/internal/utils"
)

type TodoRepository struct {
	baseUrl string
}

const Endpoint = "todos"

func NewTodoRepository(baseUrl string) *TodoRepository {
	return &TodoRepository{
		baseUrl,
	}
}

// fetch one
func (r *TodoRepository) FindById(id int) (*model.Todo, error) {
	// build URL
	url := r.baseUrl + Endpoint + "/" + strconv.Itoa(id)

	// call API util to make web request
	response, err := utils.MakeRequest(url)
	if err != nil {
		return nil, err
	}

	// unmarshall response
	var todo *model.Todo
	if err := json.Unmarshal(response, &todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// fetch many
func (r *TodoRepository) Find(numResults int) ([]*model.Todo, error) {
	// build URL
	url := r.baseUrl + Endpoint

	// call api util
	response, err := utils.MakeRequest(url)
	if err != nil {
		return nil, err
	}

	// unmarshal results into slice
	var results []*model.Todo
	if err := json.Unmarshal(response, &results); err != nil {
		return nil, err
	}

	// segment results based on requested number
	// 0 will return all results
	if numResults > 0 {
		results = results[0:numResults]
	}

	return results, nil
}
