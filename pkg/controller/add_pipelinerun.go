package controller

import (
	"github.com/bigkevmcd/task-statuses/pkg/controller/pipelinerun"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, pipelinerun.Add)
}
