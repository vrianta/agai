package session

func (h *collection) Len() int { return len(*h) }
func (h *collection) Less(i, j int) bool {
	return (*h)[i].ExpirationTime.Before((*h)[j].ExpirationTime)
}
func (h *collection) Swap(i, j int) { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *collection) Push(x any) {
	*h = append(*h, x.(*Instance))
}

// pop first element from the heap
// and return it, also remove it from the heap
func (h *collection) Pop() any {
	if len(*h) == 0 {
		return nil

	}
	*h = (*h)[1:]
	return h
}
