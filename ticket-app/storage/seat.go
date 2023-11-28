package storage

import (
	"gorm.io/gorm"
)

// Seat-related structs and methods
type Seat struct {
	ID       uint   `gorm:"primaryKey"`
	Status   string `gorm:"not null;default:'OPEN'"`
}

// CreateSeat creates a new seat
func (s *Storage) CreateSeat(seat *Seat) error {
	return s.db.Create(seat).Error
}

// GetSeatByID retrieves a seat by ID
func (s *Storage) GetSeatByID(id uint) (*Seat, error) {
	seat := &Seat{}
	err := s.db.First(seat, id).Error
	if err != nil {
		return nil, err
	}
	return seat, nil
}

// UpdateSeat updates an existing seat
func (s *Storage) UpdateSeat(seat *Seat) error {
	return s.db.Save(seat).Error
}

// DeleteSeat deletes a seat
func (s *Storage) DeleteSeat(id uint) error {
	return s.db.Delete(&Seat{}, id).Error
}

// GetSeats retrieves all seats from the database
func (s *Storage) GetSeats() ([]Seat, error) {
	var seats []Seat
	result := s.db.Find(&seats)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return seats, nil
}
