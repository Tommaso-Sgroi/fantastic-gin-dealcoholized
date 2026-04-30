/*
Package domain contains the core business entities, interfaces and rules of the application.
It has no dependency on any framework or infrastructure concern (no Gin, no SQL, no HTTP).

Domain errors are defined here so that service implementations and handlers share the same
sentinel values without importing each other.
*/
package name

// Name represents a Name entry in the system.
type Name struct {
	ID    int64
	Value string
}
