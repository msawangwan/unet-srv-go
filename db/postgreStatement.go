package db

// command
const (
	CMD_SELECT      = "SELECT"
	CMD_INSERT_INTO = "INSERT INTO"
)

// table schema
const (
	SCHEMA_PROFILE = "profile"
)

// patterns
const (
	PATTERN_STAR = "* FROM"
)

// prepared statements
const (
	STATEMENT_SELECT_ALL_PROFILES   = CMD_SELECT + " " + PATTERN_STAR + " " + SCHEMA_PROFILE
	STATEMENT_INSERT_CREATE_PROFILE = CMD_INSERT_INTO + " " + SCHEMA_PROFILE + " (profile_name, profile_uuid, hashed_gamestate, date_created, timeof_lastsave) VALUES ($1, $2, $3, $4, $5)"
)
