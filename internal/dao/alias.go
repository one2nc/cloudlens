package dao

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/config"
	"github.com/one2nc/cloudlens/internal/render"
)

var _ Accessor = (*Alias)(nil)

// Alias tracks standard and custom command aliases.
type Alias struct {
	*config.Aliases
}

// NewAlias returns a new set of aliases.
func NewAlias() *Alias {
	a := Alias{Aliases: config.NewAliases()}

	return &a
}

// Check verifies an alias is defined for this command.
func (a *Alias) Check(cmd string) bool {
	_, ok := a.Aliases.Get(cmd)
	return ok
}

// List returns a collection of aliases.
func (a *Alias) List(ctx context.Context) ([]Object, error) {
	aa, ok := ctx.Value(internal.KeyAliases).(*Alias)
	if !ok {
		return nil, fmt.Errorf("expecting *Alias but got %T", ctx.Value(internal.KeyAliases))
	}
	m := aa.ShortNames()
	oo := make([]Object, 0, len(m))
	for res, aliases := range m {
		sort.StringSlice(aliases).Sort()
		oo = append(oo, render.AliasRes{Resource: res, Aliases: aliases})
	}

	return oo, nil
}

// AsResource returns a matching resource if it exists.
func (a *Alias) AsResource(cmd string) (string, bool) {
	res, ok := a.Aliases.Get(cmd)
	if ok {
		return res, true
	}
	return "", false
}

// Get fetch a resource.
func (a *Alias) Get(_ context.Context, _ string) (Object, error) {
	return nil, errors.New("NYI!!")
}

// Ensure makes sure alias are loaded.
func (a *Alias) Ensure() (config.Alias, error) {

	return a.Alias, a.load()
}

func (a *Alias) load() error {
	return a.Load()
}
