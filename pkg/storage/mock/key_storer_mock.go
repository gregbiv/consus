package mock

import "github.com/gregbiv/consus/pkg/model"

// Storer represents a key Storer mock
type Storer struct {
	StoreError error
}

// Store is responsible for INSERT operation on key struct
func (d *Storer) Store(key *model.Key) error {
	return d.StoreError
}
