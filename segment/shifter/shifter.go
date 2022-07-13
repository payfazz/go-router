// Package shifter provide simple routing by dividing path into its segment
package shifter

// Shifter hold state of shifting segment in the path
type Shifter struct {
	tag  map[string]int
	list []string
	next int
}

// New return new Shifter with the key
func New(list []string) *Shifter {
	return &Shifter{
		tag:  nil,
		list: list,
		next: 0,
	}
}

// Reset the shifter
func (s *Shifter) Reset() {
	s.tag = nil
	s.next = 0
}

// SetNext set the index for the next Shift.
func (s *Shifter) SetNext(next int) {
	if next < 0 {
		next = 0
	}
	if size := s.Size(); next > size {
		next = size
	}
	s.next = next
}

// Shift to next segment, also telling if already in last segment
func (s *Shifter) Shift() (string, bool) {
	if s.End() {
		return "", true
	}
	ret := s.list[s.next]
	s.next++
	return ret, s.End()
}

// Unshift do reverse of Shift
func (s *Shifter) Unshift() {
	if s.next == 0 {
		return
	}
	s.next--
}

// Get i-th segment
func (s *Shifter) Get(i int) string {
	if i < 0 || i >= s.Size() {
		return ""
	}
	return s.list[i]
}

// GetRelative is same with Get, but relative to current segment
func (s *Shifter) GetRelative(d int) string {
	return s.Get(s.CurrentIndex() + d)
}

// Size return the size of segment in path
func (s *Shifter) Size() int {
	return len(s.list)
}

// CurrentIndex of shifter state
func (s *Shifter) CurrentIndex() int {
	return s.next - 1
}

// End indicated end segment in the path
func (s *Shifter) End() bool {
	return s.next == s.Size()
}

// Split return processed segment and rest of them
func (s *Shifter) Split() (done []string, rest []string) {
	done = make([]string, s.next)
	rest = make([]string, s.Size()-s.next)
	copy(done, s.list[:s.next])
	copy(rest, s.list[s.next:])
	return done, rest
}

// Tag current segment
func (s *Shifter) Tag(tag string) {
	s.TagIndex(s.CurrentIndex(), tag)
}

// TagIndex will tag i-th segment
func (s *Shifter) TagIndex(i int, tag string) {
	if i < 0 || i >= s.Size() {
		return
	}
	if s.tag == nil {
		s.tag = make(map[string]int)
	}
	s.tag[tag] = i
}

// TagRelative is same with TagIndex, but relative to current segment
func (s *Shifter) TagRelative(d int, tag string) {
	s.TagIndex(s.CurrentIndex()+d, tag)
}

// DeleteTag delete tag
func (s *Shifter) DeleteTag(tag string) {
	delete(s.tag, tag)
}

// ClearTag clear all tags on index
func (s *Shifter) ClearTag(index int) {
	var what []string
	for k, v := range s.tag {
		if v == index {
			what = append(what, k)
		}
	}
	for _, v := range what {
		delete(s.tag, v)
	}
}

// GetByTag return tagged segment
func (s *Shifter) GetByTag(tag string) (string, bool) {
	i, ok := s.tag[tag]
	if !ok {
		return "", false
	}
	return s.list[i], true
}
