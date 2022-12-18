package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	result := make(Bi)
	inn := make(Bi)
	out := startStages(inn, stages...)

	go func() {
		defer close(inn)

		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				inn <- v
			case <-done:
				return
			}
		}
	}()

	go func() {
		defer close(result)
		for {
			select {
			case v, ok := <-out:
				if !ok {
					return
				}
				result <- v
			case <-done:
				go func() {
					for range out {
					}
				}()
				return
			}
		}
	}()

	return result
}

func startStages(in In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stage(in)
	}

	return in
}
