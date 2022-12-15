// Package enrollment facilitates the binding of an Area Controller to a Building Controller.
// The Building Controller connects to the Area Controller, which implements the EnrollmentApi, and gives it
// metadata about the Smart Core network and a certificate.
//
// This package only implements the communication between the Building Controller and the Area Controller.
// It does not interact with the database or user interfaces.
package enrollment
