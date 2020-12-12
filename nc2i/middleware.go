package nc2i

import (
	"encoding/json"
	"github.com/omecodes/bome"
	"github.com/omecodes/common/httpx"
	"github.com/omecodes/common/utils/log"
	"net/http"
	"time"
)

type middleware func(handler http.Handler) http.Handler

var logHandler = httpx.Logger("nc2i")

func visits(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		date := time.Now().UnixNano()
		ctx := r.Context()
		store := visitsInfoStore(ctx)
		if store != nil {
			agent := r.Header.Get("User-Agent")
			remoteIP := r.Header.Get("Remote-Addr")
			page := r.URL.Path

			info := map[string]interface{}{
				"date":  date,
				"agent": agent,
				"ip":    remoteIP,
				"page":  page,
			}

			data, _ := json.Marshal(info)
			if data != nil {
				err := store.Save(&bome.ListEntry{
					Index: date,
					Value: string(data),
				})
				if err != nil {
					log.Error("failed to save visit info", log.Err(err))
				}
			}
		}
	})
}
