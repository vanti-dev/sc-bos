package gateway

// tasks is a helper for tracking tasks and removing them.
type tasks map[string]func()

func (t tasks) callAll() {
	for _, f := range t {
		f()
	}
}

func (t tasks) remove(k string) {
	if f, ok := t[k]; ok {
		delete(t, k)
		f()
	}
}
