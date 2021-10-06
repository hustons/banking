package transactions

type byDate []transaction

func (r byDate) Len() int {
	return len(r)
}

func (r byDate) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r byDate) Less(i, j int) bool {
	return r[i].GetCompletedDate().Before(r[j].GetCompletedDate())
}
