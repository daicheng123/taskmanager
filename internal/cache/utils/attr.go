package utils

import (
	"time"
)

const (
	AttrExpire = "expire"
)

type Attr struct {
	Name  string
	Value interface{}
}

type Attrs []*Attr

func (as Attrs) Find(name string) interface{} {
	for _, a := range as {
		if a.Name == name {
			return a.Value
		}
	}
	return nil
}

func WithExpire(t time.Duration) *Attr {
	return &Attr{
		Name:  AttrExpire,
		Value: t,
	}
}
