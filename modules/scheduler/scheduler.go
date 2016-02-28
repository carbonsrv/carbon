package scheduler

var sched chan func() // scheduler buffer

// Scheduler
// Run the scheduler
func Run() {
	sched = make(chan func(), 64)
	for {
		fn := <-sched
		go fn()
	}
}

// Add to queue
func Add(fn func()) {
	sched <- fn
}
