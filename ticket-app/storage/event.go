package storage

import (
    "gorm.io/gorm"
)

type Event struct {
    EventID        uint   `gorm:"primaryKey"`
    EventName string `gorm:"not null"`
}

func (s *Storage) CreateEvent(event *Event) error {
    return s.db.Create(event).Error
}

func (s *Storage) GetEventByID(id uint) (*Event, error) {
    event := &Event{}
    err := s.db.First(event, id).Error
    if err != nil {
        return nil, err
    }
    return event, nil
}

func (s *Storage) UpdateEvent(event *Event) error {
    return s.db.Save(event).Error
}

func (s *Storage) DeleteEvent(id uint) error {
    return s.db.Delete(&Event{}, id).Error
}

func (s *Storage) GetEvents() ([]Event, error) {
    var events []Event
    result := s.db.Find(&events)
    if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
        return nil, result.Error
    }
    return events, nil
}
