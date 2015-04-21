package scheduler

var sched chan func() // scheduler buffer

// Scheduler
func Run() {
	sched = make(chan func(), 64)
	for {
		fn := <-sched
		go fn()
	}
}
func Add(fn func()) {
	sched <- fn
}
