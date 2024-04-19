package pkg

type Set map[string]bool

func (s Set) Add(item string) {
	s[item] = true
}

func (s Set) Remove(item string) {
	delete(s, item)
}

func (s Set) Contains(item string) bool {
	_, found := s[item]
	return found
}
