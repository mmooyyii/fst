package fst

// 一种代替map[byte]*Edge的数据结构，主要将map改为有序，以便在fuzzySearch时按字典序返回，性能也许会有些提高

type ByteEdgePair struct {
	byte byte
	edge *Edge
}

type ByteMap interface {
	setEdge(key byte, val *Edge)
	setNode(key byte, val *Node)
	addOutput(key byte, val int)

	getEdge(key byte) (*Edge, bool)
	forloop() []ByteEdgePair
	upgrade() ByteMap
}

func NewByteMap() ByteMap {
	return &smallByteMap{key: make([]byte, 0), edge: make([]*Edge, 0)}
}

type smallByteMap struct {
	key  []byte
	edge []*Edge
}

func (s *smallByteMap) setEdge(key byte, val *Edge) {
	idx := s.bisectLeft(key)
	if idx == len(s.key) {
		s.key = append(s.key, key)
		s.edge = append(s.edge, val)
	} else if s.key[idx] != key {
		s.key = append(s.key, 0)
		copy(s.key[idx+1:], s.key[idx:])
		s.key[idx] = key
		s.edge = append(s.edge, nil)
		copy(s.edge[idx+1:], s.edge[idx:])
		s.edge[idx] = val
	} else {
		s.edge[idx] = val
	}
}

func (s *smallByteMap) setNode(key byte, val *Node) {
	if edge, ok := s.getEdge(key); !ok {
		panic("keyError: " + string(key))
	} else {
		edge.node = val
	}
}

func (s *smallByteMap) addOutput(key byte, val int) {
	if edge, ok := s.getEdge(key); !ok {
		panic("keyError: " + string(key))
	} else {
		edge.output += val
	}
}

func (s *smallByteMap) getEdge(key byte) (*Edge, bool) {
	idx := s.bisectLeft(key)
	if idx == len(s.key) {
		return nil, false
	} else if s.key[idx] != key {
		return nil, false
	} else {
		return s.edge[idx], true
	}
}

func (s *smallByteMap) forloop() []ByteEdgePair {
	ans := make([]ByteEdgePair, len(s.key))
	for i, key := range s.key {
		ans[i] = ByteEdgePair{
			byte: key,
			edge: s.edge[i],
		}
	}
	return ans
}

func (s *smallByteMap) bisectLeft(key byte) int {
	var mid int
	l, r := 0, len(s.key)
	for l < r {
		mid = (l + r) / 2
		if s.key[mid] < key {
			l = mid + 1
		} else {
			r = mid
		}
	}
	return l
}

const UpgradeThreshold = 100 // 暂时不知道选多少合适， 用压测确定
func (s *smallByteMap) upgrade() ByteMap {
	if len(s.key) < UpgradeThreshold {
		return s
	}
	ss := &bigByteMap{}
	for _, v := range s.forloop() {
		ss.setEdge(v.byte, v.edge)
	}
	return ss
}

type bigByteMap struct {
	data [256]*Edge
}

func (b *bigByteMap) setEdge(key byte, val *Edge) {
	b.data[key] = val
}

func (b *bigByteMap) setNode(key byte, val *Node) {
	b.data[key].node = val
}

func (b *bigByteMap) addOutput(key byte, val int) {
	b.data[key].output += val
}

func (b *bigByteMap) getEdge(key byte) (*Edge, bool) {
	edge := b.data[key]
	if edge == nil {
		return nil, false
	}
	return edge, true
}

func (b *bigByteMap) forloop() []ByteEdgePair {
	ans := make([]ByteEdgePair, 0)
	for key, val := range b.data {
		if val != nil {
			ans = append(ans, ByteEdgePair{byte: byte(key), edge: val})
		}
	}
	return ans
}

func (b *bigByteMap) upgrade() ByteMap {
	// 永不升级
	return b
}
