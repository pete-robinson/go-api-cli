package main

import (
	"flag"

	"github.com/pete-robinson/go-cli-test/internal/domain/repository"
	"github.com/pete-robinson/go-cli-test/internal/handlers/cli"
	"github.com/pete-robinson/go-cli-test/internal/services"
	log "github.com/sirupsen/logrus"
)

const baseUrl = "https://jsonplaceholder.typicode.com/"

func main() {
	// init flags
	opts := &cli.Options{}
	flag.StringVar(&opts.Command, "command", "fetchall", "Command to run")
	flag.IntVar(&opts.NumResults, "count", 10, "Number of results to return")
	flag.IntVar(&opts.Id, "id", 0, "Number of results to return")
	flag.StringVar(&opts.List, "list", "", "Fetch all results from a given list")
	flag.Parse()

	// spit out some generic logging of our flags
	log.WithFields(log.Fields{
		"command":    opts.Command,
		"numResults": opts.NumResults,
		"id":         opts.Id,
	}).Info("Booted cli")

	// boot repository
	repo := repository.NewTodoRepository(baseUrl)

	// boot todo service
	todoService := services.NewTodoService(repo)

	// init the handler
	handler := cli.New(
		cli.WithService(todoService),
		cli.WithOptions(opts),
	)

	// dispatch the command
	err := handler.Dispatch()
	if err != nil {
		log.WithField("command", opts.Command).Error(err)
	}

}
