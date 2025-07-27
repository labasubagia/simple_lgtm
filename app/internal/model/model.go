package model

import "fmt"

type DataItem struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func (v *DataItem) Validate() error {
	if v.ID == "" {
		return fmt.Errorf("ID is required")
	}
	if v.Value == "" {
		return fmt.Errorf("Value is required")
	}
	return nil
}
