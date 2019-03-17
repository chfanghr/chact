package chact

type next struct {
	ch *chainActions
}

func (n next) placeHolder() {} //used to implement NextAction interface
func (n next) JumpTo(t Tag) {
	n.ch.jumpDest = t
}
func (n next) Then(t Task) NextAction {
	if t == nil {
		panic("nil function")
	}
	n.ch.actions = append(n.ch.actions, t)
	return n
}
func (n next) Catch(e ErrorHandler) NextAction {
	if e == nil {
		panic("nil function")
	}
	n.ch.actions = append(n.ch.actions, e)
	return n
}
func (n next) Tag(t Tag) NextAction {
	if t == "" {
		panic("empty tag")
	}
	if idx := len(n.ch.actions) - 1; idx > 0 {
		n.ch.tags[t] = len(n.ch.actions) - 1
	} else {
		n.ch.tags[t] = 0
	}
	return n
}
