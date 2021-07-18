package game

import (
	"encoding/json"
	"fmt"
)

type cellType byte
type position struct {
	x, y int
}
type cell struct {
	cellType cellType
	position position
}

func (c cell) Positions() (int, int) {
	return c.position.x, c.position.y
}

const (
	DEAD  cellType = '-'
	ALIVE cellType = '+'
)

func (ct cellType) IsValid() error {
	switch ct {
	case DEAD, ALIVE:
		return nil
	}
	return fmt.Errorf("invalid cell type: %v", ct)
}

func (ct cellType) ValidOrNil() (*cellType, error) {
	if ct.String() == "" {
		return nil, nil
	}
	err := ct.IsValid()
	if err != nil {
		return &ct, err
	}
	return &ct, nil
}

func (ct cellType) String() string {
	return string(ct)
}

func (ct cellType) List() []cellType {
	return []cellType{DEAD, ALIVE}
}

func (ct cellType) StringList() []string {
	var s []string
	for _, v := range ct.List() {
		s = append(s, v.String())
	}
	return s
}

// UnmarshalJSON - implements Unmarshaler interface for cell
func (ct *cellType) UnmarshalJSON(data []byte) error {
	var s byte
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := cellType(s)
	if err := v.IsValid(); err != nil {
		return err
	}
	*ct = v
	return nil
}

// MarshalJSON - implements Marshaller interface for cell
func (ct *cellType) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.String())
}
