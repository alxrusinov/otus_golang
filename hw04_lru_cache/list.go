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

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	size int
	head *ListItem
	tail *ListItem
}

func (li *ListItem) Clear() {
	li.Next = nil
	li.Prev = nil
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	if l.head == nil {
		newHead := NewListItem(v)
		l.head = newHead
		l.tail = newHead
		l.size++
	} else {
		newHead := NewListItem(v)
		newHead.Next = l.head
		l.head.Prev = newHead
		l.head = newHead
		l.size++
	}

	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.head == nil {
		newTail := NewListItem(v)
		l.head = newTail
		l.tail = newTail
		l.size++
	} else {
		newTail := NewListItem(v)
		l.tail.Next = newTail
		newTail.Prev = l.tail
		l.tail = newTail
		l.size++
	}
	return l.tail
}

func (l *list) Remove(i *ListItem) {
	switch i {
	case l.head:
		if l.head == l.tail {
			l.head = nil
			l.tail = nil
			i.Clear()
			l.size = 0
		} else {
			l.head = l.head.Next
			l.head.Prev = nil
			i.Clear()
			l.size--
		}
	case l.tail:
		l.tail = l.tail.Prev
		l.tail.Next = nil
		i.Clear()
		l.size--
	default:
		if i.Next == nil && i.Prev == nil {
			return
		}
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		i.Clear()
		l.size--
	}
}

func (l *list) MoveToFront(i *ListItem) {
	switch i {
	case l.head:
		return
	case l.tail:
		i.Prev.Next = nil
		l.tail = i.Prev
		i.Prev = nil
		l.head.Prev = i
		i.Next = l.head
		l.head = i
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		l.head.Prev = i
		i.Prev = nil
		i.Next = l.head
		l.head = i
	}
}

func NewList() List {
	return new(list)
}

func NewListItem(v interface{}) *ListItem {
	return &ListItem{Value: v}
}
