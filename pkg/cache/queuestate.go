package cache

type States map[string]State

type State struct {
	queue *Queue
}
