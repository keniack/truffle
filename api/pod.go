package api

func (p *Pod) Exists() bool {
	_, exists := PodsMap[p.PodName]
	return exists
}
