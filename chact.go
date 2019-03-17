package chact

import (
	"context"
	"fmt"
	"sync"
)

type chainActions struct {
	mu      *sync.Mutex
	actions []action
	tags    map[Tag]int

	jumpDest Tag
	ctx      context.Context
}

func (c *chainActions) placeHolder()   {} //used to implement Chain interface
func (c *chainActions) resetJumpDest() { c.jumpDest = "" }
func (c *chainActions) Execute(p Parameters) (res Results, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.resetJumpDest()

	//insert error handler
	c.actions = append(c.actions, defaultErrorHandler)

	for i := 0; i < len(c.actions); {
		action := c.actions[i]
		if task, ok := action.(Task); ok {
			chRes, chErr := make(chan Results), make(chan error)
			go func() { res, err = task(p); chRes <- res; chErr <- err }()
			select {
			case res = <-chRes:
				err = <-chErr
				break
			case <-c.ctx.Done():
				return nil, fmt.Errorf("cancled by context")
			}

			for j := i; err != nil; { //handle the problem reported by task
				//when all the problems are resolved,exit the loop
				//find the nearest error handler
				for idx, action := range c.actions[j:] {
					if errHandler, ok := action.(ErrorHandler); ok {
						cont := false                                   //make a continue flag
						if res, err, cont = errHandler(p, err); !cont { //when continue flag is set to false
							return //just return
						} else {
							j = idx + 1 //if error handler throw another error,j will be used to navigate next error handler
							p = Parameters(res)
							break
						}
					}
				}
			}

			if c.jumpDest != "" { //handle jump call
				if idx, ok := c.tags[c.jumpDest]; !ok {
					panic("invalid tag to jump to")
				} else {
					c.jumpDest = ""
					i = idx
				}
			} else {
				i++
			}

			//every thing is fine,prepare for the next task
			p = Parameters(res)
		} else {
			i++
		}
	}
	return
}
func (c *chainActions) New(wrapper func(NextAction, Utils)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.resetJumpDest()
	c.actions = make([]action, 0)
	wrapper(next{c}, next{c})
}
func (c *chainActions) Append(wrapper func(NextAction, Utils)) {
	c.resetJumpDest()
	wrapper(next{c}, next{c})
}
func (c *chainActions) SetContext(ctx context.Context) {
	if ctx == nil {
		panic("nil context")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ctx = ctx
}
