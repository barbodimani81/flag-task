package main

// import (
// 	"testing"
// )

// func resetTestDB(t *testing.T) {
// 	t.Helper()
// 	tables := []string{"audit_logs", "dependencies", "feature_flags"}
// 	for _, tbl := range tables {
// 		_, err := testDB.Exec("DELETE FROM " + tbl)
// 		if err != nil {
// 			t.Fatalf("failed to reset table %s: %v", tbl, err)
// 		}
// 	}
// }

// func TestCreateFeatureFlag(t *testing.T) {
// 	resetTestDB(t)

// 	err := createFeatureFlag(testDB, "auth", nil)
// 	if err != nil {
// 		t.Fatalf("failed to create flag: %v", err)
// 	}

// 	var count int
// 	err = testDB.QueryRow(`SELECT COUNT(*) FROM feature_flags WHERE name = "auth"`).Scan(&count)
// 	if err != nil || count != 1 {
// 		t.Fatalf("expected 1 auth flag, got %d, err: %v", count, err)
// 	}
// }
