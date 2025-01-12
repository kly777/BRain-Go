package relation

import (
	"database/sql"
	"time"
)

type Relation struct {
	ID        int64     `json:"id"`
	Relata    int64     `json:"relata"`   // Card ID
	Relation  int64     `json:"relation"` // Card ID
	Position  int       `json:"position"`
	Describe  string    `json:"describe"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func Create(db *sql.DB, relata, relation int64, position int, describe string) (*Relation, error) {
	query := `
		INSERT INTO relations (relata, relation, position, describe, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := db.Exec(query, relata, relation, position, describe, now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Relation{
		ID:        id,
		Relata:    relata,
		Relation:  relation,
		Position:  position,
		Describe:  describe,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func GetByID(db *sql.DB, id int64) (*Relation, error) {
	query := `
		SELECT id, relata, relation, position, describe, created_at, updated_at
		FROM relations
		WHERE id = ?
	`

	row := db.QueryRow(query, id)

	var r Relation
	err := row.Scan(&r.ID, &r.Relata, &r.Relation, &r.Position, &r.Describe, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func Update(db *sql.DB, id int64, relata, relation int64, position int, describe string) error {
	query := `
		UPDATE relations
		SET relata = ?, relation = ?, position = ?, describe = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := db.Exec(query, relata, relation, position, describe, time.Now(), id)
	return err
}

func Delete(db *sql.DB, id int64) error {
	query := `
		DELETE FROM relations
		WHERE id = ?
	`

	_, err := db.Exec(query, id)
	return err
}

func ListByRelata(db *sql.DB, relata int64) ([]*Relation, error) {
	query := `
		SELECT id, relata, relation, position, describe, created_at, updated_at
		FROM relations
		WHERE relata = ?
		ORDER BY position ASC
	`

	rows, err := db.Query(query, relata)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relations []*Relation
	for rows.Next() {
		var r Relation
		err := rows.Scan(&r.ID, &r.Relata, &r.Relation, &r.Position, &r.Describe, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
		relations = append(relations, &r)
	}

	return relations, nil
}
