package option

import (
	. "github.com/FreeSamples00/gonads"
)

func Map[I any, O any](o Option[I], fn func(Option[I]) O) Option[O] {
	if o.IsSome() {
		return Some(fn(o))
	}
	return o
}
