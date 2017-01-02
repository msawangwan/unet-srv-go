package db

// redis keys (those ending with ":" are to be concatonated with unique value)
const (
	// string
	KEY_IS_LOADED_IN_MEMORY = "world:is_loaded:"
	KEY_SESSION             = "session:"

	// sets
	KEY_NAMES_TAKEN       = "profile_name:taken"
	KEY_SESSION_AVAILABLE = "session:active"

	// hash
	KEY_WORLD_NODES      = "world:valid_nodes:"
	KEY_SESSION_INSTANCE = "session:active:instance:"
)

// redis default and placeholder values
const (
	VAL_INIT = "_init_"
)
