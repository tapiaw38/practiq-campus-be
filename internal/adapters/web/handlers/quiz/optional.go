package quiz

import (
	"bytes"
	"encoding/json"
)

// optionalString tells "the client did not mention this field" apart from
// "the client set it to null".
//
// Both decode into a nil *string, which made an update that omitted a field
// indistinguishable from one clearing it — so saving a quiz's title through a
// form that had no date inputs silently erased its schedule. Absent now means
// leave it alone.
type optionalString struct {
	Present bool
	Value   *string
}

func (o *optionalString) UnmarshalJSON(data []byte) error {
	o.Present = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		o.Value = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	o.Value = &s
	return nil
}
