package mock

import "github.com/gregbiv/sandbox/pkg/model"

// Updater represents a key Updater mock
type Updater struct {
	UpdateError error
}

// Update is responsible for UPDATE operation on key struct
func (d *Updater) Update(key *model.Key) error {
	return d.UpdateError
}
