package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cockroachdb/apd/v3"
)

type Decimal struct {
	*apd.Decimal
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	if d.Decimal == nil {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", d.Decimal)), nil
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		d.Decimal = nil
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		var f float64
		if err := json.Unmarshal(data, &f); err != nil {
			return err
		}
		d.Decimal = apd.New(0, 0)
		_, err = d.Decimal.SetFloat64(f)
		return err
	}

	var dec apd.Decimal
	if _, _, err := dec.SetString(s); err != nil {
		return err
	}
	d.Decimal = &dec
	return nil
}

func (d Decimal) Value() (driver.Value, error) {
	if d.Decimal == nil {
		return nil, nil
	}
	return d.Decimal.String(), nil
}

func (d *Decimal) Scan(value interface{}) error {
	if value == nil {
		d.Decimal = nil
		return nil
	}

	var str string
	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return errors.New("tipo incompatível para Decimal.Scan")
	}

	var dec apd.Decimal
	if _, _, err := dec.SetString(str); err != nil {
		return err
	}
	d.Decimal = &dec
	return nil
}
