package nc2i

import (
	"fmt"
	"github.com/omecodes/bome"
	"github.com/omecodes/common/httpx"
	"github.com/omecodes/common/utils/log"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	appRoute      = "/app"
	metricsRoute  = "/metrics"
	messagesRoute = "/messages"
)

func serveWebApp(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, appRoute)

	if !strings.HasPrefix(filename, "/") {
		filename = "/" + filename
	}

	if filepath.Ext(filename) == "" || filename == "/" {
		filename = path.Join(filename, "index.html")
	}

	ctx := r.Context()

	mime, size, content, err := getResourceFile(ctx, filename)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	defer func() {
		if err != nil {
			log.Error("close file reader", log.Err(err))
		}
	}()

	buf := make([]byte, 2048)
	done := false
	w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
	w.Header().Set("Content-Type", mime)

	for !done {
		n, err := content.Read(buf)
		if err != nil {
			done = err == io.EOF
			if !done {
				log.Error("failed to send content", log.Err(err))
				return
			}
		}
		_, _ = w.Write(buf[:n])
	}
}

func saveMessage(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" && contentType != "text/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.ContentLength > 2048 {
		log.Error("message size exceeds limit of 2048 bytes")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("could get message from request body", log.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	msgStore := messages(ctx)

	date := time.Now().UnixNano() / 1e6

	err = msgStore.Save(&bome.ListEntry{
		Index: date,
		Value: string(data),
	})
	if err != nil {
		log.Error("failed to save message", log.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
