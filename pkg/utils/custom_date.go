package utils

import (
	"database/sql/driver"
	"errors"
	"strings"
	"time"
)

// CustomDate untuk parsing tanggal dalam format "YYYY-MM-DD"
type CustomDate struct {
	time.Time
}

// UnmarshalJSON → parsing dari JSON ke time.Time
func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		cd.Time = time.Time{}
		return nil
	}

	// Parsing format tanggal sederhana
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	cd.Time = t
	return nil
}

// MarshalJSON → konversi ke string "YYYY-MM-DD"
func (cd CustomDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + cd.Format("2006-01-02") + `"`), nil
}

// String → representasi string-nya
func (cd CustomDate) String() string {
	return cd.Format("2006-01-02")
}

// ✅ Implement database/sql interfaces biar GORM bisa save & read
func (cd CustomDate) Value() (driver.Value, error) {
	if cd.IsZero() {
		return nil, nil
	}
	return cd.Time, nil // GORM akan simpan sebagai timestamp
}

func (cd *CustomDate) Scan(value interface{}) error {
	if value == nil {
		cd.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		cd.Time = v
	case []byte:
		t, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		cd.Time = t
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		cd.Time = t
	default:
		return errors.New("unsupported type for CustomDate")
	}
	return nil
}
