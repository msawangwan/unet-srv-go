package db

// redis keys (those ending with ":" are to be concatonated with unique value)
const (
	// string
	KEY_IS_LOADED_IN_MEMORY = "world:is_loaded:"

	// sets
	KEY_NAMES_TAKEN = "profile_name:taken"

	// hash
	KEY_WORLD_NODES = "world:valid_nodes:"
)

// redis default and placeholder values
const (
	VAL_INIT = "_init_"
)
