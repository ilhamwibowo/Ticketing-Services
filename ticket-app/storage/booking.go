package storage

import (
	"gorm.io/gorm"
)

type Booking struct {
	ID          uint   `gorm:"primaryKey"`
	SeatID      uint   `gorm:"not null"`
	InvoiceID   string `gorm:"not null"`
	PaymentURL  string `gorm:"not null"`
	Status      string `gorm:"not null;default:'ONGOING'"`
}

// CreateBooking creates a new booking
func (s *Storage) CreateBooking(booking *Booking) error {
	return s.db.Create(booking).Error
}

// UpdateBookingStatus updates the status of a booking
func (s *Storage) UpdateBookingStatus(id uint, status string) error {
	return s.db.Model(&Booking{}).Where("id = ?", id).Update("status", status).Error
}

func (s *Storage) UpdateBookingStatusByInvoiceID(invoiceID, status string) error {
	return s.db.Model(&Booking{}).Where("invoice_id = ?", invoiceID).Update("status", status).Error
}

// GetBookingByInvoiceID retrieves a booking by its invoice ID
func (s *Storage) GetBookingByInvoiceID(invoiceID string) (*Booking, error) {
	booking := &Booking{}
	err := s.db.Where("invoice_id = ?", invoiceID).First(booking).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil, nil if the record is not found
		}
		return nil, err
	}
	return booking, nil
}