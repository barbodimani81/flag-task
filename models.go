package main

import "time"

type FeatureFlag struct {
	ID        int
	Name      string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Dependency struct {
	ID          int
	FlagID      int
	DependsOnID int
}

type AuditLog struct {
	ID        int
	FlagID    int
	Action    string
	Reason    string
	Actor     string
	Timestamp time.Time
}
