package model

type Mergeable interface {
	Merge(src *Mergeable)
}
