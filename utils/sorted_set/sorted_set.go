package sorted_set

import (
	"fmt"
	"github.com/aruyuna9531/skiplist"
)

type SortedSet[Key comparable, Value skiplist.ISkiplistElement[Key]] struct {
	s    *skiplist.SkipList[Key]
	dict map[Key]Value
}

func NewSortedSet[Key comparable, Value skiplist.ISkiplistElement[Key]]() *SortedSet[Key, Value] {
	return &SortedSet[Key, Value]{
		s:    skiplist.NewSkipList[Key](),
		dict: make(map[Key]Value),
	}
}

func (s *SortedSet[Key, Value]) Add(v Value) error {
	if _, exist := s.dict[v.Key()]; exist {
		return fmt.Errorf("SortedSet Add error: key %v exist", v.Key())
	}
	err := s.s.Add(v)
	if err != nil {
		return err
	}
	s.dict[v.Key()] = v
	return nil
}

func (s *SortedSet[Key, Value]) Delete(k Key) error {
	err := s.s.DeleteByKey(k)
	if err != nil {
		return err
	}
	delete(s.dict, k)
	return nil
}

func (s *SortedSet[Key, Value]) Update(v Value) (orVal Value, err error) {
	// 非原子操作注意
	orVal = s.GetByKey(v.Key())
	err = s.s.DeleteByKey(v.Key())
	if err != nil {
		return
	}
	err = s.s.Add(v)
	s.dict[v.Key()] = v
	return
}

func (s *SortedSet[Key, Value]) GetByKey(k Key) (v Value) {
	v = s.dict[k]
	return
}

func (s *SortedSet[Key, Value]) GetByRank(r int32) (v Value, err error) {
	elem, err := s.s.GetElementByRank(r)
	if err != nil {
		return
	}
	v = s.dict[elem.Key()]
	return
}

func (s *SortedSet[Key, Value]) GetRankByKey(k Key) (rank int32) {
	rank, _ = s.s.GetRankByKey(k)
	return
}

func (s *SortedSet[Key, Value]) RemoveByRank(rank int32) (originElement Value, err error) {
	// 非原子操作注意
	i, err := s.s.GetElementByRank(rank)
	if err != nil {
		return
	}
	originElement = s.dict[i.Key()]
	err = s.s.DeleteByKey(i.Key())
	if err != nil {
		return
	}
	delete(s.dict, i.Key())
	return
}

func (s *SortedSet[Key, Value]) GetCount() int32 {
	return s.s.GetElementsCount()
}

func (s *SortedSet[Key, Value]) GetRange(start, end int32) (res []Value, err error) {
	ret, err := s.s.GetRange(start, end)
	if err != nil {
		return
	}
	for _, v := range ret {
		res = append(res, s.dict[v.Key()])
	}
	return
}

func (s *SortedSet[Key, Value]) GetAll() []Value {
	ret := make([]Value, 0)
	for _, v := range s.dict {
		ret = append(ret, v)
	}
	return ret
}
