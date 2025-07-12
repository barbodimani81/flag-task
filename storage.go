package main

import (
	"database/sql"
	"fmt"
)

func createFeatureFlag(db *sql.DB, name string, deps []string) error {
	// Start TX
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert flag
	res, err := tx.Exec(`INSERT INTO feature_flags (name) VALUES (?)`, name)
	if err != nil {
		return err
	}
	flagID, _ := res.LastInsertId()

	// Insert dependencies
	for _, depName := range deps {
		var depID int
		err := tx.QueryRow(`SELECT id FROM feature_flags WHERE name = ?`, depName).Scan(&depID)
		if err != nil {
			return fmt.Errorf("dependency %s not found", depName)
		}

		if causesCycle(tx, int(flagID), depID) {
			return fmt.Errorf("circular dependency detected with %s", depName)
		}

		_, err = tx.Exec(`INSERT INTO dependencies (flag_id, depends_on_id) VALUES (?, ?)`, flagID, depID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func toggleFlag(db *sql.DB, name string, enable bool, actor, reason string) error {
	// Get flag
	var id int
	err := db.QueryRow(`SELECT id FROM feature_flags WHERE name = ?`, name).Scan(&id)
	if err != nil {
		return err
	}

	if enable {
		// Check dependencies
		missing, err := getMissingDeps(db, id)
		if err != nil {
			return err
		}
		if len(missing) > 0 {
			return fmt.Errorf("missing active dependencies: %v", missing)
		}
	} else {
		// Auto-disable dependents
		disableDependents(db, id, actor, "auto-disabled due to parent flag being disabled")
	}

	_, err = db.Exec(`UPDATE feature_flags SET enabled = ? WHERE id = ?`, enable, id)
	if err != nil {
		return err
	}

	// Log
	return insertAuditLog(db, id, "toggle", actor, reason)
}

func getMissingDeps(db *sql.DB, flagID int) ([]string, error) {
	rows, err := db.Query(`
		SELECT f.name FROM dependencies d
		JOIN feature_flags f ON f.id = d.depends_on_id
		WHERE d.flag_id = ? AND f.enabled = FALSE
	`, flagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missing []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		missing = append(missing, name)
	}
	return missing, nil
}

func disableDependents(db *sql.DB, flagID int, actor, reason string) error {
	rows, err := db.Query(`SELECT flag_id FROM dependencies WHERE depends_on_id = ?`, flagID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var dependents []int
	for rows.Next() {
		var depID int
		rows.Scan(&depID)
		dependents = append(dependents, depID)
	}

	for _, id := range dependents {
		_, err := db.Exec(`UPDATE feature_flags SET enabled = FALSE WHERE id = ?`, id)
		if err != nil {
			return err
		}
		insertAuditLog(db, id, "auto-disable", actor, reason)
		disableDependents(db, id, actor, reason) // Recursive
	}
	return nil
}
