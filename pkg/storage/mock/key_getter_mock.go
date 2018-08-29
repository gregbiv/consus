package mock

import "github.com/gregbiv/consus/pkg/model"

// Getter represents a key Getter mock
type Getter struct {
	Key      *model.Key
	KeyError error
}

// GetAll gets all keys
func (g *Getter) GetAll(filterStr string, activeOnly bool) ([]*model.Key, error) {
	if g.Key == nil {
		return []*model.Key{}, g.KeyError
	}

	return []*model.Key{g.Key}, g.KeyError
}

// GetByID gets a key by ID
func (g *Getter) GetByID(ID string, activeOnly bool) (*model.Key, error) {
	return g.Key, g.KeyError
}
