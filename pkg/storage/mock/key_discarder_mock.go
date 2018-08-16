package mock

// Discarder represents a key Discarder mock
type Discarder struct {
	DiscardError error
}

// Discard is responsible for DELETE operation on key struct
func (d *Discarder) Discard(ID string) error {
	return d.DiscardError
}

// Truncate is responsible for TRUNCATE operation on key struct
func (d *Discarder) Truncate() error {
	return d.DiscardError
}
