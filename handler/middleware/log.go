package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Logging struct {
	Timestamp time.Time
	Latency   int64
	Path      string
	OS        string
}

func Log(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()

		defer func(started time.Time) {
			logging := Logging{
				Timestamp: started,
				Latency:   int64(time.Since(started)),
				Path:      r.URL.Path,
				OS:        GetOSCtx(r.Context()),
			}

			j, _ := json.Marshal(logging)

			fmt.Println(string(j))

		}(started)

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
