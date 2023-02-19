package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		localOut := make(Bi)
		localIn := stage(out)

		go func() {
			defer close(localOut)

			for {
				select {
				case <-done:
					return
				case v, ok := <-localIn:
					if !ok {
						return
					}

					select {
					case <-done:
						return
					default:
						localOut <- v
					}
				}
			}
		}()
		out = localOut
	}

	return out
}
