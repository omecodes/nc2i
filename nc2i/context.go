package nc2i

import (
	"context"
	"github.com/omecodes/bome"
	"net/http"
)

type ctxMessagesStore struct{}
type ctxVisitsInfo struct{}
type ctxResFS struct{}
type ctxResDir struct{}
type ctxDataDir struct{}

func contextWithMessages(parent context.Context, list *bome.JSONList) context.Context {
	return context.WithValue(parent, ctxMessagesStore{}, list)
}

func messages(ctx context.Context) *bome.JSONList {
	o := ctx.Value(ctxMessagesStore{})
	if o == nil {
		return nil
	}
	return o.(*bome.JSONList)
}

func contextWithVisitsInfoStore(parent context.Context, list *bome.JSONList) context.Context {
	return context.WithValue(parent, ctxVisitsInfo{}, list)
}

func visitsInfoStore(ctx context.Context) *bome.JSONList {
	o := ctx.Value(ctxVisitsInfo{})
	if o == nil {
		return nil
	}
	return o.(*bome.JSONList)
}

func contextWithResFS(parent context.Context, fs http.FileSystem) context.Context {
	return context.WithValue(parent, ctxResFS{}, fs)
}

func resFS(ctx context.Context) http.FileSystem {
	o := ctx.Value(ctxResFS{})
	if o == nil {
		return nil
	}
	return o.(http.FileSystem)
}

func contextWithExternalResDir(parent context.Context, dirname string) context.Context {
	return context.WithValue(parent, ctxResDir{}, dirname)
}

func externalResDir(ctx context.Context) string {
	o := ctx.Value(ctxResDir{})
	if o == nil {
		return ""
	}
	return o.(string)
}

func contextWithDataDir(parent context.Context, dirname string) context.Context {
	return context.WithValue(parent, ctxDataDir{}, dirname)
}

func appDataDir(ctx context.Context) string {
	o := ctx.Value(ctxDataDir{})
	if o == nil {
		return ""
	}
	return o.(string)
}
