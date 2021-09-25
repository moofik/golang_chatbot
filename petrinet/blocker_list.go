package petrinet

type BlockerList struct {
	index    int
	blockers []*Blocker
}

func (l *BlockerList) add(b *Blocker) {
	l.blockers = append(l.blockers, b)
}

func (l *BlockerList) next() (*Blocker, bool) {
	has := false

	if l.index < len(l.blockers) {
		has = true
	}

	if has {
		blocker := l.blockers[l.index]
		l.index++
		return blocker, true
	}

	return nil, false
}

func (l *BlockerList) current() *Blocker {
	if l.index != 0 {
		return l.blockers[l.index-1]
	}

	return nil
}

func (l *BlockerList) has(code string) bool {
	for _, blocker := range l.blockers {
		if code == blocker.code {
			return true
		}
	}

	return false
}

func (l *BlockerList) empty() bool {
	return len(l.blockers) == 0
}

func (l *BlockerList) count() int {
	return len(l.blockers)
}
