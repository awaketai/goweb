package logs

import "sync"

type logger struct {
	opt *options
	mu sync.Mutex
	entryPool *sync.Pool
}

// func NewLogger(opts ...Option) *logger {
// 	logger := &logger{
// 		opt:       initOptions(opts...),
// 		mu:        sync.Mutex{},
// 		entryPool: &sync.Pool{
// 			New: func() any {
// 				return entr
// 			},
// 		},
// 	}
// }