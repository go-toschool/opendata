package postgres

import (
	"github.com/Finciero/opendata/sagittarius"
)

// SagittariusService implements aiolos.Sagittarius interface with postgres as DB.
type SagittariusService struct {
	Store SQLExecutor
}

// Create ...
func (ss *SagittariusService) Create(c *aiolos.Callback) error {
	return nil
}

// Read ...
func (ss *SagittariusService) Read(id string) (*aiolos.Callback, error) {
	return nil, nil
}

// Update ...
func (ss *SagittariusService) Update(c *aiolos.Callback) error {
	return nil
}

// Delete ...
func (ss *SagittariusService) Delete(c *aiolos.Callback) error {
	return nil
}
