package nc2i

import (
	"encoding/json"
	"github.com/omecodes/bome"
	"github.com/omecodes/common/httpx"
	"github.com/omecodes/common/utils/log"
	"net/http"
	"path"
	"time"
)

var visitsExcludedExtensions = []string{
	".js",
	".css",
	".ttf",
	".otf",
	".json",
	".csh",
	".svg",
	".png",
	".ico",
}
var visitsExcludedBaseNames = []string{
	"NOTICES",
	"robot.txt",
}
var visitsExcludedPaths = []string{
	"/",
}

type middleware func(handler http.Handler) http.Handler

var logHandler = httpx.Logger("nc2i")

func visits(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)

		agent := r.Header.Get("User-Agent")
		remoteIP := r.Header.Get("Remote-Addr")
		page := r.URL.Path

		for _, x := range visitsExcludedPaths {
			if r.URL.Path == x {
				return
			}
		}

		extension := path.Ext(r.URL.Path)
		for _, x := range visitsExcludedExtensions {
			if extension == x {
				return
			}
		}

		base := path.Base(r.URL.Path)
		for _, x := range visitsExcludedBaseNames {
			if base == x {
				return
			}
		}

		date := time.Now().UnixNano()
		ctx := r.Context()
		store := visitsInfoStore(ctx)
		if store != nil {

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
