package linkedList

import (
	"fmt"
)

type EqualFunc[T any] func(a T, b T) bool

type ILinkedList[T any] interface {
	AddFront(item T)
	AddBack(item T)
	Remove(item T, isEqual EqualFunc[T]) error
	RemoveByIndex(index int) T
	Find(element T, isEqual EqualFunc[T]) (T, error)
	FindByIndex(index int) (T, error)
	GetSize() int
	PrintList()
	Clear()
}

type node[T any] struct {
	element  T
	nextNode *node[T]
}

type linkedList[T any] struct {
	head   *node[T]
	tail   *node[T]
	length int
}

func New[T any]() *linkedList[T] {
	return &linkedList[T]{
		head:   nil,
		tail:   nil,
		length: 0,
	}
}

func (list *linkedList[T]) GetSize() int {
	return list.length
}

func (list *linkedList[T]) AddFront(item T) {
	newNode := &node[T]{element: item, nextNode: nil}
	if list.head == nil {
		list.head = newNode
		list.tail = newNode
		list.length++
		return
	}

	newNode.nextNode = list.head
	list.head = newNode
	list.length++
}

func (list *linkedList[T]) AddBack(item T) {
	newNode := &node[T]{element: item, nextNode: nil}
	if list.head == nil {
		list.head = newNode
		list.tail = newNode
		list.length++
		return
	}

	list.tail.nextNode = newNode
	list.tail = newNode
	list.length++
}

func (list *linkedList[T]) Find(element T, isEqual EqualFunc[T]) (T, error) {
	current := list.head

	for current != nil {
		if isEqual(element, current.element) {
			return current.element, nil
		}
		current = current.nextNode
	}

	var defaultValue T
	return defaultValue, ErrElementNotFound
}

func (list *linkedList[T]) FindByIndex(index int) (T, error) {
	var defaultValue T
	if index > list.length || index < 0 {
		return defaultValue, ErrIndexOutOfRange
	}

	current := list.head
	var iterator = 0
	for current != nil {
		if index == iterator {
			return current.element, nil
		}
		current = current.nextNode
		iterator++
	}

	return defaultValue, ErrElementNotFound
}

func (list *linkedList[T]) Remove(element T, isEqual EqualFunc[T]) error {
	if list.length == 0 {
		return ErrorElementNotRemoved
	}

	if isEqual(element, list.head.element) {
		list.head = list.head.nextNode
		list.length--
		if list.head == nil {
			list.tail = nil
		}
		return nil
	}

	var prev *node[T]
	current := list.head
	for current != nil {
		if isEqual(current.element, element) {
			prev.nextNode = current.nextNode
			list.length--

			if current.nextNode == nil {
				list.tail = prev
			}
			return nil
		}

		prev = current
		current = current.nextNode
	}

	return ErrorElementNotRemoved
}

func (list *linkedList[T]) RemoveByIndex(index int) T {
	var element T
	if index == 0 {
		element = list.head.element
		list.head = list.head.nextNode
		list.length--

		if list.head == nil {
			list.tail = nil
		}

		return element
	}

	var iterator int
	var prev *node[T]
	current := list.head
	for current != nil {
		if index == iterator {
			prev.nextNode = current.nextNode
			list.length--

			if current.nextNode == nil {
				list.tail = prev
			}
			return current.element
		}

		prev = current
		current = current.nextNode
		iterator++
	}

	return element
}

func (list *linkedList[T]) PrintList() {
	current := list.head
	for current != nil {
		fmt.Printf("%v ", current.element)
		current = current.nextNode
	}
	fmt.Printf("\n")
}

func (list *linkedList[T]) Clear() {
	list.head = nil
	list.tail = nil
	list.length = 0
}
