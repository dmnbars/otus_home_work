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
	count int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	current := &ListItem{Value: v}
	l.pushFront(current)

	return current
}

func (l *list) PushBack(v interface{}) *ListItem {
	current := &ListItem{Value: v}
	l.count++

	if l.back != nil {
		current.Prev = l.back
		current.Prev.Next = current
		l.back = current

		return current
	}

	l.front = current
	l.back = current

	return current
}

func (l *list) Remove(i *ListItem) {
	l.count--

	if l.count == 0 {
		l.front = nil
		l.back = nil

		return
	}

	if i.Prev != nil && i.Next != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev

		return
	}

	if i.Prev != nil {
		i.Prev.Next = nil
		l.back = i.Prev
		if l.count == 1 {
			l.front = l.back
		}
	}
	if i.Next != nil {
		i.Next.Prev = nil
		l.front = i.Next
		if l.count == 1 {
			l.back = l.front
		}
	}
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.pushFront(i)
}

func (l *list) pushFront(i *ListItem) {
	l.count++
	if l.front != nil {
		i.Next = l.front
		i.Next.Prev = i
		l.front = i

		return
	}

	l.front = i
	l.back = i
}

func NewList() List {
	return &list{}
}
