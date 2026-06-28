// Package workflow is the workflow/task engine that backs every multi-step
// process in the platform.
//
// A workflow is a DAG of Tasks. Each Task has a type (e.g. "wirepay.transfer"),
// inputs, and an SLA. Tasks are scheduled onto RabbitMQ; the engine advances
// them as their dependencies complete.
package workflow

import (
	"context"

	"github.com/scape-labs/monorepo/shared/idgen"
)

// Workflow is the persisted form of a workflow definition.
//
// TODO: replace with a real schema (jsonb in Postgres, versioned).
type Workflow struct {
	ID    idgen.ID
	Name  string
	Tasks []Task
}

// Task is one node in a workflow DAG.
//
// TODO: add Status, Attempts, LastError, etc.
type Task struct {
	ID     idgen.ID
	Type   string
	Inputs map[string]any
}

// Engine runs workflows.
//
// TODO: implement — schedule tasks, advance on completion, retry on failure.
type Engine struct {
	// TODO: fields.
}

// Start kicks off a workflow instance.
//
// TODO: implement.
func (e *Engine) Start(ctx context.Context, wf Workflow) error {
	_ = ctx
	_ = wf
	// TODO: implement.
	return nil
}

// Advance is called when a task finishes and we need to schedule its
// downstream successors.
//
// TODO: implement.
func (e *Engine) Advance(ctx context.Context, taskID idgen.ID) error {
	_ = ctx
	_ = taskID
	// TODO: implement.
	return nil
}
