package queue

type States map[string]State

type State struct {
	queue *Queue
}
