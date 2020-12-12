package nc2i

import (
	"encoding/json"
	"github.com/mssola/user_agent"
	"github.com/omecodes/bome"
	"github.com/omecodes/common/httpx"
	"github.com/omecodes/common/utils/log"
	"net/http"
	"time"
)

var logHandler = httpx.Logger("nc2i")

func middleware(final http.Handler) http.Handler {
	return logHandler.Handle(visits(final))
}

func visits(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)

		ctx := r.Context()
		store := metrics(ctx)
		if store != nil {
			agent := user_agent.New(r.Header.Get("User-Agent"))
			remoteIP := r.Header.Get("Remote-Address")
			page := r.URL.Path

			info := map[string]interface{}{
				"date":  time.Now().UnixNano() / 1e6,
				"agent": agent,
				"ip":    remoteIP,
				"page":  page,
			}

			data, _ := json.Marshal(info)
			if data != nil {
				err := store.Save(&bome.ListEntry{
					Index: time.Now().UnixNano(),
					Value: string(data),
				})
				if err != nil {
					log.Error("failed to save visit info", log.Err(err))
				}
			}
		}
	})
}
