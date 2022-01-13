package cli

import (
	"errors"

	"github.com/pete-robinson/go-cli-test/internal/domain/model"
	log "github.com/sirupsen/logrus"
)

const (
	CMD_FETCHALL  string = "fetchall"
	CMD_FETCHONE         = "fetchone"
	CMD_FETCHLIST        = "fetchlist"
)

// accepted cli options
type Options struct {
	Command    string
	NumResults int
	Id         int
	List       string
}

// struct for this handler
type CliHandler struct {
	options *Options
	service Service
}

// the signature for the service we're hoping to receive
type Service interface {
	FetchAll(int) ([]*model.Todo, error)
	FetchOne(int) (*model.Todo, error)
	FetchList(string) ([]*model.Todo, error)
}

// constructor
func New(options ...func(*CliHandler)) *CliHandler {
	handler := &CliHandler{}
	for _, o := range options {
		o(handler)
	}

	return handler
}

// dispatch command
func (s *CliHandler) Dispatch() error {
	log.Infof("command: %s", s.options.Command)

	switch s.options.Command {
	// fetch all entries
	case CMD_FETCHALL:
		res, err := s.service.FetchAll(s.options.NumResults)
		if err != nil {
			return err
		}

		// loop over results and spit them out - add a hacky time delay
		for _, r := range res {
			s.LogTodoOutput(r)
		}

	// fetch a single entry
	case CMD_FETCHONE:
		res, err := s.service.FetchOne(s.options.Id)
		if err != nil {
			return err
		}

		// spit result out in logging
		s.LogTodoOutput(res)
	case CMD_FETCHLIST:
		res, err := s.service.FetchList(s.options.List)
		if err != nil {
			return err
		}

		for _, r := range res {
			s.LogTodoOutput(r)
		}
	default:
		return errors.New("Invalid command requested")
	}

	return nil
}

// log out a data set
func (s *CliHandler) LogTodoOutput(todo *model.Todo) {
	log.WithFields(log.Fields{
		"id":        todo.Id,
		"userId":    todo.UserId,
		"title":     todo.Title,
		"completed": todo.Completed,
	}).Info("Response received")
}

// functional options
func WithService(svc Service) func(*CliHandler) {
	return func(c *CliHandler) {
		c.service = svc
	}
}

func WithOptions(options *Options) func(*CliHandler) {
	return func(c *CliHandler) {
		c.options = options
	}
}
