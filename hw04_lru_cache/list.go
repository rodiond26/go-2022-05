package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

// Node
type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

// Open doubly linked list
type list struct {
	length int
	head   *ListItem
	tail   *ListItem
}

func NewList() List {
	return &list{}
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newHead := &ListItem{Value: v}
	switch {

	case l.length == 0:
		l.head = newHead
		l.tail = newHead

	case l.length == 1:
		l.head = newHead
		l.head.Next = l.tail
		l.tail.Prev = l.head

	case l.length > 1:
		newHead.Next = l.head
		l.head.Prev = newHead
		l.head = newHead
	}

	l.length++
	return newHead
}

func (l *list) PushBack(v interface{}) *ListItem {
	newTail := &ListItem{Value: v}
	switch {

	case l.length == 0:
		l.head = newTail
		l.tail = newTail

	case l.length == 1:
		l.tail = newTail
		l.tail.Prev = l.head
		l.head.Next = l.tail

	case l.length > 1:
		newTail.Prev = l.tail
		l.tail.Next = newTail
		l.tail = newTail
	}

	l.length++
	return newTail
}

func (l *list) Remove(rem *ListItem) {
	switch {

	case l.length == 1:
		l.head = nil
		l.tail = nil

	case l.length == 2 && rem == l.head:
		l.head = l.tail
		l.head.Next = nil
		l.head.Prev = nil

	case l.length == 2 && rem == l.tail:
		l.tail = l.head
		l.tail.Next = nil
		l.tail.Prev = nil

	case l.length > 2 && rem == l.head:
		l.head = l.head.Next
		l.head.Prev = nil

	case l.length > 2 && rem == l.tail:
		l.tail = l.tail.Prev
		l.tail.Next = nil

	default:
		rem.Next.Prev = rem.Prev
		rem.Prev.Next = rem.Next
	}

	l.length--
}

func (l *list) MoveToFront(newHead *ListItem) {
	switch {

	case l.head == newHead:

	case l.length == 1:

	case l.length == 2 && l.tail == newHead:
		l.head, l.tail = l.tail, l.head
		l.head.Next = l.tail
		l.head.Prev = nil
		l.tail.Prev = l.head
		l.tail.Next = nil

	case l.length > 2 && l.tail == newHead:
		l.tail = l.tail.Prev
		l.tail.Next = nil
		l.head.Prev = newHead
		newHead.Next = l.head
		l.head = newHead
		l.head.Prev = nil

	default:
		newHead.Prev.Next = newHead.Next
		newHead.Next.Prev = newHead.Prev
		l.head.Prev = newHead
		newHead.Next = l.head
		newHead.Prev = nil
		l.head = newHead
	}
}
