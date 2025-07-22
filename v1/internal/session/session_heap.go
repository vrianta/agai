package session

func (h *collection) Len() int { return len(*h) }
func (h *collection) Less(i, j int) bool {
	return (*h)[i].ExpirationTime.Before((*h)[j].ExpirationTime)
}
func (h *collection) Swap(i, j int) { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *collection) Push(x any) {
	*h = append(*h, x.(*Instance))
}

func (h *collection) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}
