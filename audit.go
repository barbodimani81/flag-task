package main

import "database/sql"

func insertAuditLog(db *sql.DB, flagID int, action, actor, reason string) error {
	_, err := db.Exec(`INSERT INTO audit_logs (flag_id, action, actor, reason) VALUES (?, ?, ?, ?)`,
		flagID, action, actor, reason)
	return err
}
