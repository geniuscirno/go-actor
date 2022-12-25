package actor

import "go-actor/core"

type SpawnOption func(opts *core.SpawnOptions)

func Name(name string) SpawnOption {
	return func(opts *core.SpawnOptions) {
		opts.Name = name
	}
}
