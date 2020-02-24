package fsm

import (
	"context"

	"github.com/filecoin-project/go-statestore"
)

// EventName is the name of an event
type EventName interface{}

// Context provides access to the statemachine inside of a state handler
type Context interface {
	// Context returns the golang context for this context
	Context() context.Context

	// Event initiates a state transition with the named event.
	//
	// The call takes a variable number of arguments that will be passed to the
	// callback, if defined.
	//
	// It will return nil if the event is one of these errors:
	//
	// - event X does not exist
	//
	// - arguments don't match expected transition
	Event(event EventName, args ...interface{}) error
}

// ApplyTransitionFunc modifies the state further in addition
// to modifying the state key. It the signature
// func applyTransition<StateType, T extends any[]>(s stateType, args ...T)
// and then an event can be dispatched on context or group
// with the form .Event(Name, args ...T)
type ApplyTransitionFunc interface{}

// StateKeyField is the name of a field in a state struct that serves as the key
// by which the current state is identified
type StateKeyField string

// StateKey is a value for the field in the state that represents its key
// that uniquely identifies the state
// in practice it must have the same type as the field in the state struct that
// is designated the state key and must be comparable
type StateKey interface{}

// TransitionMap is a map from src state to destination state
type TransitionMap map[StateKey]StateKey

// Events describes the different events that can happen in a state machine
type Events []EventBuilder

// StateType is a type for a state, represented by an empty concrete value for a state
type StateType interface{}

// Environment are externals dependencies will be needed by this particular state machine
type Environment interface{}

// StateHandler is called upon entering a state after
// all events are processed. It should have the signature
// func stateHandler<StateType, Environment>(ctx Context, environment Environment, state StateType) error
type StateHandler interface{}

// StateHandlers is a map between states and their handlers
type StateHandlers map[StateKey]StateHandler

// Notifier should be a function that takes two parameters,
// a native event type and a statetype
// -- nil means no notification
// it is called after every successful state transition
// with the even that triggered it
type Notifier func(eventName EventName, state StateType)

// Group is a manager of a group of states that follows finite state machine logic
type Group interface {

	// Begin initiates tracking with a specific value for a given identifier
	Begin(id interface{}, userState interface{}) error

	// Send sends the given event name and parameters to the state specified by id
	// it will error if there are underlying state store errors or if the parameters
	// do not match what is expected for the event name
	Send(id interface{}, name EventName, args ...interface{}) (err error)

	// SendSync will block until the given event is actually processed, and
	// will return an error if the transition was not possible given the current
	// state
	SendSync(ctx context.Context, id interface{}, name EventName, args ...interface{}) (err error)

	// Get gets state for a single state machine
	Get(id interface{}) *statestore.StoredState

	// List outputs states of all state machines in this group
	// out: *[]StateT
	List(out interface{}) error

	// Stop stops all state machines in this group
	Stop(ctx context.Context) error
}

// Parameters are the parameters that define a finite state machine
type Parameters struct {
	// required

	// Environment is the environment in which the state handlers operate --
	// used to connect to outside dependencies
	Environment Environment

	// StateType is the type of state being tracked. Should be a zero value of the state struct, in
	// non-pointer form
	StateType StateType

	// StateKeyField is the field in the state struct that will be used to uniquely identify the current state
	StateKeyField StateKeyField

	// Events is the list of events that that can be dispatched to the state machine to initiate transitions.
	// See EventDesc for event properties
	Events []EventBuilder

	// StateHandlers - functions that will get called each time the machine enters a particular
	// state. this is a map of state key -> handler.
	StateHandlers StateHandlers

	// optional

	// Notifier is a function that gets called on every successful event processing
	// with the event name and the new state
	Notifier Notifier
}
