package ui

import (
	"context"
	"time"

	"github.com/one2nc/cloud-lens/internal/dao"
	"github.com/one2nc/cloud-lens/internal/model"
	"github.com/one2nc/cloud-lens/internal/render"
)

// Lister represents a viewable resource.
type Lister interface {
	// Get returns a resource instance.
	Get(ctx context.Context, path string) (dao.Object, error)
}

// Tabular represents a tabular model.
type Tabular interface {
	Lister

	// Empty returns true if model has no data.
	Empty() bool

	// Count returns the model data count.
	Count() int

	// Peek returns current model data.
	Peek() *render.TableData

	// Watch watches a given resource for changes.
	//Watch(context.Context) error

	// Refresh forces a new refresh.
	Refresh(context.Context) error

	// SetRefreshRate sets the model watch loop rate.
	SetRefreshRate(time.Duration)

	// AddListener registers a model listener.
	AddListener(model.TableListener)

	// RemoveListener unregister a model listener.
	RemoveListener(model.TableListener)
}
