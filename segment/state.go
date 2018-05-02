package segment

type keyType struct{}

var key keyType

type state struct {
	list []string
	next int
}

func (s *state) shift() (string, bool) {
	if s.next == len(s.list) {
		return "", false
	}
	ret := s.list[s.next]
	s.next++
	return ret, true
}

func (s *state) unshift() {
	if s.next == 0 {
		return
	}
	s.next--
}

func (s *state) get(i int) string {
	if i < 0 || i >= len(s.list) {
		return ""
	}
	return s.list[i]
}
