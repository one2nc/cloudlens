package dao

import (
	"context"
)

type Object interface{}

// Getter represents a resource getter.
type Getter interface {
	// Get return a given resource.
	Get(ctx context.Context, path string) (Object, error)
}

// Lister represents a resource lister.
type Lister interface {
	// List returns a resource collection.
	List(ctx context.Context) ([]Object, error)
}

type Accessor interface {
	Lister
	Getter
}
