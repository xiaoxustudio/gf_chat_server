package imc

import (
	"context"
	"gf_server/utility/tw"
)

type Imc struct{}

func New() *Imc {
	return &Imc{}
}

func (r *Imc) connect() int {
	tw.Tw(context.Background(), "阿松大")
	return 1
}
