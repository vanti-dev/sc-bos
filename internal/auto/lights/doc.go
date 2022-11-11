// Package lights provides automated control of lighting.
// Information is read from sensors, occupancy and/or brightness, that then informs decisions about whether lights
// should be on and at what brightness.
//
// This package provides these automatic lighting functions:
//
//  1. Lights in occupied spaces are on, lights in unoccupied spaces are off
//  2. Lights that are on are dimmer if it's bright outside
//
// # Implementation approach
//
// The core automation logic can be found in the processState function in logic.go.
// This code receives the combined knowledge of the system and works out what needs to be done to satisfy the automatic
// outcomes we desire.
// Here you'll find the code that says "if any sensor reports a state of OCCUPIED, then the lights should be on".
//
// Collecting and updating the shared knowledge is a process we call patching, inspired by version control systems.
// ReadState holds all the information we can know and it is updated by a sequence of Patcher.Patch calls. These
// Patcher instances come from the different sensors and are applied to the shared state one at a time in state.go#readStateChanges.
//
// Invoking the logic when the state changes happens in processStateChanges which takes care to avoid unnecessary work
// by skipping analysing state that has already been superseded by more state changes.
//
// # Implementation finer details
//
// To avoid duplicating commands and to provide some sense of determinism the logic function is allowed to complete
// before being called again, even if more state changes are available.
// The context passed to the func is only cancelled if the automation as a whole is stopped.
// This means that we shouldn't get into a state where the logic has turns some lights on but not others because a new
// version of the state is available.
// Actions successfully performed by the logic are recorded in WriteState which is consulted before new actions are performed.
// If the logic wants to turn the lights on, then it can check WriteState to see if the light has already been asked to
// turn on, the action can then be skipped as needed.
//
// One of the clauses for the logic is to wait a configured time after a space becomes unoccupied before turning the
// lights off. To accomplish this the logic could use time.After but that would cause the logic to block for a long time,
// during which the automation is not reacting to new information. Instead the logic returns a TTL value that indicates
// to the calling func that the analysis of the state and the subsequent actions are only correct until TTL expires, at
// which point the logic should be rerun.
package lights
