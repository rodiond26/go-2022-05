package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		in := make(Bi)
		close(in)
		return in
	}

	out := in

	for _, stage := range stages {
		out = stage(doStage(out, done))
	}

	return doStage(out, done)
}

func doStage(inputCh In, terminateCh In) Out {
	resultCh := make(Bi)

	go func() {
		defer func() {
			close(resultCh)
			for range inputCh {
			}
		}()

		for {
			select {
			case val, ok := <-inputCh:
				if !ok {
					return
				}
				select {
				case resultCh <- val:
				case <-terminateCh:
					return
				}
			case <-terminateCh:
				return
			}
		}
	}()

	return resultCh
}
