package parse

// Env represents a collection of environment variables with efficient lookup capabilities.
// It maintains environment variables in "KEY=VALUE" format and provides an indexed
// mapping for fast retrieval. Duplicate keys are handled by keeping only the first
// occurrence and marking subsequent duplicates as empty strings.
type Env struct {
	env     []string
	indexes map[string]int
}

// NewEnv creates a new Env instance from a slice of environment variable strings.
// Each string should be in the format "KEY=VALUE". The function automatically
// handles duplicate keys by keeping only the first occurrence.
//
// Example:
//
//	env := NewEnv([]string{"HOME=/home/user", "PATH=/usr/bin", "HOME=/duplicate"})
//	// The second HOME entry will be ignored
func NewEnv(env []string) *Env {
	e := &Env{env: env}
	e.init()
	return e
}

// init initializes the Env instance by building an index map for efficient lookups.
// It processes all environment strings, extracts keys, and handles duplicates by
// keeping only the first occurrence of each key.
func (e *Env) init() {
	envs := e.env
	indexes := make(map[string]int)
	for i, s := range envs {
		for j := 0; j < len(s); j++ {
			if s[j] == '=' {
				key := s[:j]
				if _, ok := indexes[key]; !ok {
					indexes[key] = i // first mention of key
				} else {
					envs[i] = ""
				}
				break
			}
		}
	}
	e.indexes = indexes
}

// Get retrieves the value of an environment variable by its key.
// It returns an empty string if the key is not found or if the
// environment variable has no value after the '=' character.
//
// Example:
//
//	value := env.Get("HOME")  // Returns "/home/user" for "HOME=/home/user"
//	missing := env.Get("MISSING")  // Returns ""
func (e *Env) Get(key string) string {
	env := e.indexes
	i, ok := env[key]
	if !ok {
		return ""
	}
	s := e.env[i]
	// get the value of the key, eg. "FOO=bar" -> "bar"
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return s[i+1:]
		}
	}
	return ""
}

// Has checks whether an environment variable with the given key exists.
// It returns true if the key is present, false otherwise.
//
// Example:
//
//	exists := env.Has("HOME")    // Returns true if HOME is set
//	missing := env.Has("MISSING") // Returns false if MISSING is not set
func (e *Env) Has(key string) bool {
	if _, ok := e.indexes[key]; ok {
		return ok
	}
	return false
}

// Set sets an environment variable with the given key and value.
// If the key already exists, it updates the value. If not, it adds a new entry.
// The method maintains the internal index for efficient future lookups.
//
// Example:
//
//	env.Set("NEW_VAR", "value")     // Adds a new environment variable
//	env.Set("HOME", "/new/home")    // Updates existing HOME variable
func (e *Env) Set(key, value string) {
	envStr := key + "=" + value

	if i, exists := e.indexes[key]; exists {
		// Key already exists, update it
		e.env[i] = envStr
	} else {
		// Key doesn't exist, add it
		e.env = append(e.env, envStr)
		e.indexes[key] = len(e.env) - 1
	}
}
