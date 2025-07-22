package session

func (h *SessionHeap) Len() int           { return len(*h) }
func (h *SessionHeap) Less(i, j int) bool { return (*h)[i].Expiry.Before((*h)[j].Expiry) }
func (h *SessionHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *SessionHeap) Push(x any) {
	*h = append(*h, x.(*Instance))
}

func (h *SessionHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}
