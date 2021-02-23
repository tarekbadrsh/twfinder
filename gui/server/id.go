// ID type definition, and unique ID generation.

package server

import (
	"strconv"
	"sync/atomic"
)

// ID is the type of the ids of the components.
type ID int

// Note: it is intentional that base 10 is used (and not e.q. 16),
// because it is handled as a number at the client side (in JavaScript).
// It has some benefit like no need to quote IDs, simpler code generation.

// Converts an ID to a string.
func (id ID) String() string {
	return strconv.Itoa(int(id))
}

// AtoID converts a string to ID.
func AtoID(s string) (ID, error) {
	id, err := strconv.Atoi(s)

	if err != nil {
		return ID(0), err
	}
	return ID(id), nil
}

// Component ID generation and provider

// Last used value for ID
var lastID int64

// nextCompID returns a unique component ID
// First ID given is 1.
func nextCompID() ID {
	return ID(atomic.AddInt64(&lastID, 1))
}
