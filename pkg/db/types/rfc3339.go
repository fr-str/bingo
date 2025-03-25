package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type RFC3339 struct {
	time.Time
}

func (t RFC3339) Scan(value any) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case string:
		return t.UnmarshalText([]byte(v))
	case []byte:
		return t.UnmarshalText(v)
	default:
		return fmt.Errorf("cannot scan %T", value)
	}
}

func (t RFC3339) Value() (driver.Value, error) {
	return t.Format(time.RFC3339), nil
}
