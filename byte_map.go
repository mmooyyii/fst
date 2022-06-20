package fst

// 一种代替map[byte]*Edge的数据结构，主要将map改为有序，以便在fuzzySearch时按字典序返回，性能也许会有些提高

type byteMap interface {
	addEdge(key byte, val *Edge)
	setNode(key byte, val *Node)
	addOutput(key byte, val int)
	getEdge(key byte) (*Edge, bool)
	forloop() map[byte]*Edge
	upgrade() byteMap
}

type pair struct {
	key  byte
	node *Edge
}

type smallByteMap struct {
	key  []byte
	edge []*Edge
}

func (s *smallByteMap) addEdge(key byte, val *Edge) {
}

func (s *smallByteMap) setNode(key byte, val *Node) {
	//TODO implement me
	panic("implement me")
}

func (s *smallByteMap) addOutput(key byte, val int) {
	//TODO implement me
	panic("implement me")
}

func (s *smallByteMap) getEdge(key byte) (*Edge, bool) {
	//TODO implement me
	panic("implement me")
}

func (s *smallByteMap) forloop() map[byte]*Edge {
	//TODO implement me
	panic("implement me")
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

func (s *smallByteMap) upgrade() byteMap {
	//TODO implement me
	panic("implement me")
}

type bigByteMap struct {
	data  [256]*Edge
	count int
}

func (b *bigByteMap) addEdge(key byte, val *Edge) {
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

func (b *bigByteMap) forloop() map[byte]*Edge {
	ans := make(map[byte]*Edge)
	for key, val := range b.data {
		if val != nil {
			ans[byte(key)] = val
		}
	}
	return ans
}

func (b *bigByteMap) upgrade() byteMap {
	// 永不升级
	return b
}
