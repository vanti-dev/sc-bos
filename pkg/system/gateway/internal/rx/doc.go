// Package rx provides reactive data structures.
// Reactive in this sense means that changes to the data structure are broadcast to listeners.
//
// Notification in this package do not block, instead each mutation method returns a channel that is closed when all listeners have been notified.
package rx
