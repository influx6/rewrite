package stackexpr

import (
	. "github.com/influx6/rewrite"
)

const noValueMessage = "unable to use Description in empty state"

// Description manages a stack of Applicable implementing
// objects which allows popping and pushing value.
type Description struct {
	err    error
	stacks []Applicable
}

// SetErr sets the error returned when Description.Apply is
// called.
func (s *Description) SetErr(err error) {
	s.err = err
}

// Err returns possible attached error of a giving Description.
func (s *Description) Err() error {
	return s.err
}

// Push adds a new item into the Applicable list.
func (s *Description) Push(item Applicable) {
	s.stacks = append(s.stacks, item)
}

// Root returns first Applicable object in stack.
// Usually the first Applicable is the source and
// root of all defined Definitions.
//
// If there are no elements in stack, a default EmptyApplicable
// is returned.
func (s *Description) Root() Applicable {
	if len(s.stacks) == 0 {
		panic(noValueMessage)
	}
	return s.stacks[0]
}

// IsUsable returns true/false if stack is empty.
func (s *Description) IsUsable() bool {
	return len(s.stacks) == 0 || s.stacks == nil
}

// Current returns current Applicable object in stack.
//
// If there are no elements in stack, a default EmptyApplicable
// is returned.
func (s *Description) Current() Applicable {
	var target = s.get()
	if target == nil {
		panic(noValueMessage)
	}
	return target
}

// Pop pops recent stack to the last used stack.
// If called iteratively then all items will be removed from stack.
//
// If there are no elements in stack, a default EmptyApplicable
// is returned.
func (s *Description) Pop() Applicable {
	if len(s.stacks) == 0 {
		panic(noValueMessage)
	}

	elem := s.stacks[len(s.stacks)-1]
	s.stacks = s.stacks[:len(s.stacks)-1]
	return elem
}

// Release will pop the current top elements on the stack
// applying it to it's parent.
func (s *Description) Release() (parent Applicable, child Applicable) {
	if !s.IsUsable() {
		panic(noValueMessage)
	}

	child = s.Pop()
	parent = s.get()
	return
}

func ApplyTo(stack Stack, definitions ...Definition) {
	for _, definition := range definitions {
		if stack.Err() != nil {
			return
		}
		definition(stack)
	}
}

// Get returns current Applicable object in stack.
func (s *Description) get() Applicable {
	if len(s.stacks) == 0 {
		return nil
	}
	return s.stacks[len(s.stacks)-1]
}

func PopApplicable() Definition {
	return func(root Stack) {
		root.Pop()
	}
}

func PushApplicable(t Applicable) Definition {
	return func(root Stack) {
		root.Push(t)
	}
}

func Describe(definitions ...Definition) DefinitionMiddleware {
	return func(source Applicable) (Applicable, error) {
		var stack Description
		stack.Push(source)
		ApplyTo(&stack, definitions...)
		return source, stack.err
	}
}
