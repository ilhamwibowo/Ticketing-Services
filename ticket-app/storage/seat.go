package storage

import (
    "gorm.io/gorm"
	"errors"
)

type Seat struct {
    ID     uint   `gorm:"primaryKey"`
    SeatNumber string `gorm:"not null"`
    EventID uint `gorm:"not null"`
    Status string `gorm:"not null;default:'OPEN'"`
}

func (s *Storage) CreateSeat(seat *Seat) error {
    return s.db.Create(seat).Error
}

func (s *Storage) GetSeatByID(id uint) (*Seat, error) {
    seat := &Seat{}
    err := s.db.First(seat, id).Error
    if err != nil {
        return nil, err
    }
    return seat, nil
}

func (s *Storage) UpdateSeat(seat *Seat) error {
    return s.db.Save(seat).Error
}

func (s *Storage) DeleteSeat(id uint) error {
    return s.db.Delete(&Seat{}, id).Error
}

func (s *Storage) GetSeatsByEventID(eventID uint) ([]Seat, error) {
    var seats []Seat
    result := s.db.Where("event_id = ?", eventID).Find(&seats)
    if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
        return nil, result.Error
    }
    return seats, nil
}

func (s *Storage) GetEmptySeatsByEventID(eventID uint) ([]Seat, error) {
    var emptySeats []Seat
    result := s.db.Where("event_id = ? AND status = ?", eventID, "OPEN").Find(&emptySeats)
    if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
        return nil, result.Error
    }
    return emptySeats, nil
}

func (s *Storage) GetSeatByEventIDAndNumber(eventID uint, seatNumber string) (*Seat, error) {
	seat := &Seat{}
	err := s.db.Where("event_id = ? AND seat_number = ?", eventID, seatNumber).First(seat).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return seat, nil
}

func (s *Storage) UpdateSeatStatusByID(seatID uint, status string) error {
	seat := &Seat{ID: seatID}
	err := s.db.Model(seat).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}


func (s *Storage) GetSeats() ([]Seat, error) {
	var seats []Seat
	result := s.db.Find(&seats)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return seats, nil
}
