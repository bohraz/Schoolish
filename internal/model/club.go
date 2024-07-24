package model

import (
	"database/sql"
	"encoding/json"
	"strconv"
)

type Club struct {
	ID           int `json:"-"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	DateCreated  string `json:"dateCreated"`
	OwnerId      sql.NullInt64 `json:"ownerId"`
	Integrations sql.NullString `json:"integrations"`
	Members      []User `json:"members"`
}

func (c *Club) UnmarshalJSON(data []byte) error {
	var tmp struct {
		ID string `json:"id"`
		Name string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	var err error
	c.ID, err = strconv.Atoi(tmp.ID)
	if err != nil {
		return err
	}
	c.Name = tmp.Name
	c.Description = tmp.Description

	return nil
}