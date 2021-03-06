package controller

import (
	"github.com/bigkevmcd/commit-status-tracker/pkg/controller/pipelinerun"
	"github.com/bigkevmcd/commit-status-tracker/pkg/controller/taskrun"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, pipelinerun.Add, taskrun.Add)
}
