package null

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// UUID is a nullable UUID.
// It will marshal to null if null.
// Blank string input will be considered null.
type UUID struct {
	UUID  uuid.UUID
	Valid bool
}

// NewUUID creates a new UUID.
func NewUUID(u uuid.UUID, valid bool) UUID {
	return UUID{
		UUID:  u,
		Valid: valid,
	}
}

// UUIDFrom creates a new UUID that will always be valid.
func UUIDFrom(u uuid.UUID) UUID {
	return NewUUID(u, true)
}

// UUIDFromPtr creates a new UUID that be null if u is nil.
func UUIDFromPtr(u *uuid.UUID) UUID {
	if u == nil {
		return NewUUID(uuid.Nil, false)
	}

	return NewUUID(*u, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (u UUID) ValueOrZero() uuid.UUID {
	if !u.Valid {
		return uuid.Nil
	}
	return u.UUID
}

// MarshalJSON implements json.Marshaler.
func (u UUID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(u.UUID)
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *UUID) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}

	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		u.UUID = uuid.MustParse(x)
	case map[string]interface{}:
		err = json.Unmarshal(data, &u.UUID)
	case nil:
		u.Valid = false
		return nil
	default:
		err = errors.New("null: couldn't unmarshal into value of type null.UUID")
	}

	u.Valid = err == nil
	return err
}

// MarshalText implements encoding.TextMarshaler.
func (u UUID) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte("null"), nil
	}

	return u.UUID.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *UUID) UnmarshalText(data []byte) error {
	id, err := uuid.ParseBytes(data)
	if err != nil {
		u.Valid = false
		return err
	}
	u.UUID = id
	u.Valid = true
	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (u UUID) MarshalBinary() ([]byte, error) {
	if !u.Valid {
		return []byte(nil), nil
	}
	return u.UUID[:], nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (u *UUID) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return fmt.Errorf("null: invalid UUID (got %d bytes)", len(data))
	}
	copy(u.UUID[:], data)
	u.Valid = true
	return nil
}

// Scan implements the Scanner interface.
func (u *UUID) Scan(value interface{}) error {
	if value == nil {
		u.UUID, u.Valid = uuid.Nil, false
		return nil
	}

	err := u.UUID.Scan(value)
	if err != nil {
		u.Valid = false
		return err
	}

	u.Valid = true
	return nil
}

// Value implements the driver Valuer interface.
func (u UUID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}

	return u.UUID.Value()
}
