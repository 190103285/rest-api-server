package main

import (
	"database/sql"
)

type coffee struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *coffee) getCoffee(db *sql.DB) error {
	return db.QueryRow("SELECT name FROM coffee WHERE id=$1", c.ID).Scan(&c.ID, &c.Name)
}

func (c *coffee) updateCoffee(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE coffee SET name=$1 WHERE id=$3",
			c.Name, c.ID)

	return err
}

func (c *coffee) deleteCoffee(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM coffee WHERE id=$1", c.ID)
	return err
}

func (c *coffee) createCoffee(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO coffee(id, name) VALUES($1, $2) RETURNING id", c.Name).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

func getCoffees(db *sql.DB, start, count int) ([]coffee, error) {
	rows, err := db.Query("SELECT id, name FROM coffee LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	coffees := []coffee{}

	for rows.Next() {
		var c coffee
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		coffees = append(coffees, c)
	}

	return coffees, nil
}
