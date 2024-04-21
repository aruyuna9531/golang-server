package templated_list

import (
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	l := List[int]{}
	for i := 0; i < 100; i++ {
		l.PushBack(i)
	}

	for k := l.Front(); k != &l.root; k = k.next {
		fmt.Printf("%d ", k.Value)
	}
	fmt.Printf("len: %d", l.Len())

	//firstPtr := l.Front() // point at index 35
	//for k := 0; k < 35; k++ {
	//	firstPtr = firstPtr.next
	//}

	lastPtr := l.Front() // point at index 70
	for k := 0; k < 35; k++ {
		lastPtr = lastPtr.next
	}
	fmt.Printf("len: %d\n", l.Len())

	l.removeRange(l.Front(), nil)

	for k := l.Front(); k != &l.root; k = k.next {
		fmt.Printf("%d ", k.Value)
	}
	fmt.Printf("len: %d\n", l.Len())
}

func BenchmarkList(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New[int]()
		l.PushBackList()
	}
}
