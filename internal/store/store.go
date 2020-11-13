package store

import "github.com/lutomas/go-project-stub/types"

type Store interface {
	// Methods to manipulate DB
	CreateAbc(in *types.Abc) (*types.Abc, error)
	GetAbc(in *types.Abc) (*types.Abc, error)
	ListAbc(in *types.Abc) (*types.Abc, error)
	UpdateAbc(in *types.Abc) (*types.Abc, error)
}
