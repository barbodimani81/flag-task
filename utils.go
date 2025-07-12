package main

import (
	"database/sql"
)

func causesCycle(tx *sql.Tx, newFlagID int, depID int) bool {
	// Check if newFlagID is reachable from depID via DFS
	visited := make(map[int]bool)
	return dfsCycle(tx, depID, newFlagID, visited)
}

func dfsCycle(tx *sql.Tx, currentID int, targetID int, visited map[int]bool) bool {
	if currentID == targetID {
		return true
	}
	if visited[currentID] {
		return false
	}
	visited[currentID] = true

	rows, err := tx.Query(`SELECT depends_on_id FROM dependencies WHERE flag_id = ?`, currentID)
	if err != nil {
		return false // assume no cycle if DB error (safe fallback)
	}
	defer rows.Close()

	for rows.Next() {
		var next int
		rows.Scan(&next)
		if dfsCycle(tx, next, targetID, visited) {
			return true
		}
	}
	return false
}
