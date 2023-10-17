package flip

import "fmt"

type Idempotent struct {
	counter int
}

func (i Idempotent) GetIdempotent() string {
	i.counter++
	return fmt.Sprintf("%v", i.counter)
}
