// Package rx provides reactive data structures.
// Reactive in this sense means that changes to the data structure are broadcast to listeners.
//
// Writes block until all listeners have received the change, as such listeners should be quick.
package rx
