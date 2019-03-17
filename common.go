package chact

import (
	"context"
	"sync"
)

//exported types
type (
	EmptyInterfaces []interface{}
	Parameters      EmptyInterfaces
	Results         EmptyInterfaces

	Task         func(Parameters) (Results, error)
	ErrorHandler func(Parameters, error) (Results, error, bool)
	Tag          string

	Chain interface {
		Execute(Parameters) (Results, error)
		New(func(NextAction, Utils))
		Append(func(NextAction, Utils))
		SetContext(context.Context)
		placeHolder()
	}
	NextAction interface {
		Then(Task) NextAction
		Catch(ErrorHandler) NextAction
		Tag(Tag) NextAction
		placeHolder()
	}
	Utils interface {
		JumpTo(Tag)
	}
)

//internal types
type (
	action interface{} //can be a Task or a ErrorHandler
)

var defaultErrorHandler ErrorHandler = func(p Parameters, e error) (Results, error, bool) {
	return Results(p), e, false //return the original result and error from the function which generate the error and a no continue execution flag
}

func AsParameters(ifs ...interface{}) Parameters {
	return Parameters(EmptyInterfaces(ifs))
}
func AsResults(ifs ...interface{}) Results {
	return Results(EmptyInterfaces(ifs))
}

func NewChain(ctx context.Context) Chain {
	if ctx == nil {
		ctx = context.Background()
	}
	return &chainActions{mu: new(sync.Mutex), tags: make(map[Tag]int), ctx: ctx}
}
