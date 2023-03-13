package parse

type Env struct {
	env     []string
	indexes map[string]int
}

func NewEnv(env []string) *Env {
	e := &Env{env: env}
	e.init()
	return e
}

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

func (e *Env) Get(key string) string {
	env := e.indexes
	i, ok := env[key]
	if !ok {
		return ""
	}
	s := e.env[i]
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return s[i+1:]
		}
	}
	return ""
}

func (e *Env) Has(key string) bool {
	if _, ok := e.indexes[key]; ok {
		return ok
	}
	return false
}
