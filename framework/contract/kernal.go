package contract

import "net/http"

const KernelKey = "web:kernel"

type Kernel interface {
	HttpEngine() http.Handler
}
