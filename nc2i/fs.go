package nc2i

import (
	"context"
	"github.com/omecodes/common/errors"
	"io"
	"os"
	"path"
	"path/filepath"
)

var supportedMimes = map[string]string{
	".js":   "text/javascript",
	".html": "text/html",
	".csh":  "text/x-script.csh",
	".css":  "text/css",
	".svg":  "image/svg+xml",
}

func getResourceFile(ctx context.Context, filename string) (string, int64, io.ReadCloser, error) {
	fs := resFS(ctx)
	if fs == nil {
		return getExternalResourceFile(ctx, filename)
	}

	file, err := fs.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return getExternalResourceFile(ctx, filename)
		}
		return "", 0, nil, err
	}

	stats, err := file.Stat()
	if err != nil {
		return "", 0, nil, err
	}

	return supportedMimes[path.Ext(filename)], stats.Size(), file, nil
}

func getExternalResourceFile(ctx context.Context, filename string) (string, int64, io.ReadCloser, error) {
	resDir := externalResDir(ctx)
	if resDir != "" {
		return "", 0, nil, errors.NotFound
	}

	fileFullName := filepath.Join(resDir, filename)

	file, err := os.Open(fileFullName)
	if err != nil {
		return "", 0, nil, err
	}

	stats, err := file.Stat()
	if err != nil {
		return "", 0, nil, err
	}

	return supportedMimes[path.Ext(filename)], stats.Size(), file, nil
}
