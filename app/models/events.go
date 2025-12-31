package models

import (
	"time"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// JSONB type for PostgreSQL JSONB fields
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

// TimeOnly represents a time-only value (HH:MM:SS) for PostgreSQL TIME type
type TimeOnly struct {
	time.Time
}

// Scan implements the sql.Scanner interface for TimeOnly
func (t *TimeOnly) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v
		return nil
	case string:
		// Parse time string in format "HH:MM:SS" or "HH:MM"
		layouts := []string{"15:04:05", "15:04", time.RFC3339}
		for _, layout := range layouts {
			if parsed, err := time.Parse(layout, v); err == nil {
				// Use today's date with the parsed time
				now := time.Now()
				t.Time = time.Date(now.Year(), now.Month(), now.Day(),
					parsed.Hour(), parsed.Minute(), parsed.Second(), 0, time.UTC)
				return nil
			}
		}
		return fmt.Errorf("cannot scan %v into TimeOnly", value)
	case []byte:
		return t.Scan(string(v))
	default:
		return fmt.Errorf("cannot scan %v into TimeOnly", value)
	}
}

// Value implements the driver.Valuer interface for TimeOnly
func (t TimeOnly) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	// Return time in HH:MM:SS format
	return t.Time.Format("15:04:05"), nil
}

// MarshalJSON implements json.Marshaler for TimeOnly
func (t TimeOnly) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	// Return time in HH:MM format for JSON
	return []byte(fmt.Sprintf(`"%s"`, t.Time.Format("15:04"))), nil
}

// UnmarshalJSON implements json.Unmarshaler for TimeOnly
func (t *TimeOnly) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.Time = time.Time{}
		return nil
	}
	// Remove quotes
	str := string(data)
	if len(str) > 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	// Parse time string
	layouts := []string{"15:04:05", "15:04", time.RFC3339}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, str); err == nil {
			now := time.Now()
			t.Time = time.Date(now.Year(), now.Month(), now.Day(),
				parsed.Hour(), parsed.Minute(), parsed.Second(), 0, time.UTC)
			return nil
		}
	}
	return fmt.Errorf("cannot unmarshal %s into TimeOnly", string(data))
}

type EventType struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `json:"name"`
}

type EventCategory struct {
	ID              uint                `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string              `json:"name"`
	EventTypeID     uint                `json:"event_type_id"`
	EventType       EventType           `gorm:"foreignKey:EventTypeID" json:"event_type,omitempty"`
	SubCategories   []EventSubCategory  `gorm:"foreignKey:EventCategoryID" json:"sub_categories,omitempty"`
}

type EventSubCategory struct {
	ID              uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string        `json:"name"`
	EventCategoryID uint          `json:"event_category_id"`
	EventCategory   EventCategory `gorm:"foreignKey:EventCategoryID" json:"event_category,omitempty"`
	Description     string        `json:"description,omitempty"`
	CreatedOn       time.Time     `gorm:"autoCreateTime" json:"created_on,omitempty"`
	UpdatedOn       *time.Time    `gorm:"autoUpdateTime" json:"updated_on,omitempty"`
}

type EventDetails struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	EventTypeID uint      `json:"event_type_id"`
	EventType   EventType `gorm:"foreignKey:EventTypeID" json:"event_type,omitempty"`

	EventCategoryID uint          `json:"event_category_id"`
	EventCategory   EventCategory `gorm:"foreignKey:EventCategoryID" json:"event_category,omitempty"`

	EventSubCategoryID *uint            `json:"event_sub_category_id,omitempty"`
	EventSubCategory   *EventSubCategory `gorm:"foreignKey:EventSubCategoryID" json:"event_sub_category,omitempty"`

	Scale           string     `json:"scale,omitempty"`
	Theme           string     `json:"theme,omitempty"`
	StartDate       time.Time  `json:"start_date,omitempty"`
	EndDate         time.Time  `json:"end_date,omitempty"`
	DailyStartTime  *TimeOnly  `gorm:"type:time" json:"daily_start_time,omitempty"`
	DailyEndTime    *TimeOnly  `gorm:"type:time" json:"daily_end_time,omitempty"`
	SpiritualOrator string     `json:"spiritual_orator,omitempty"`
	Language        string     `json:"language,omitempty"`

	Country    string `json:"country,omitempty"`
	State      string `json:"state,omitempty"`
	City       string `json:"city,omitempty"`
	District   string `json:"district,omitempty"`
	PostOffice string `json:"post_office,omitempty"`
	Pincode    string `json:"pincode,omitempty"`
	Address    string `json:"address,omitempty"`
	AddressType string `json:"address_type,omitempty"`
	PoliceStation string `json:"police_station,omitempty"`
	AreaCovered  string `json:"area_covered,omitempty"`

	BeneficiaryMen   int `json:"beneficiary_men"`
	BeneficiaryWomen int `json:"beneficiary_women"`
	BeneficiaryChild int `json:"beneficiary_child"`
	InitiationMen    int `json:"initiation_men"`
	InitiationWomen  int `json:"initiation_women"`
	InitiationChild  int `json:"initiation_child"`

	// Branch association (nullable - optional field for backward compatibility)
	BranchID *uint   `json:"branch_id,omitempty"`
	Branch   *Branch `gorm:"foreignKey:BranchID" json:"branch,omitempty"`

	Status string `gorm:"default:'incomplete';type:varchar(20)" json:"status,omitempty"`

	CreatedOn time.Time  `json:"created_on,omitempty"`
	UpdatedOn *time.Time `json:"updated_on,omitempty"`
	CreatedBy string     `json:"created_by,omitempty"`
	UpdatedBy string     `json:"updated_by,omitempty"`

	// Note: Draft fields removed - now using separate event_drafts table
}
