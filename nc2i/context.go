package nc2i

import (
	"context"
	"github.com/omecodes/bome"
	"net/http"
)

type ctxMessagesStore struct{}
type ctxMailerSourceName struct{}
type ctxNotificationEmail struct{}
type ctxVisitsInfo struct{}
type ctxResFS struct{}
type ctxResDir struct{}
type ctxDataDir struct{}

func contextWithMessages(parent context.Context, list *bome.JSONList) context.Context {
	return context.WithValue(parent, ctxMessagesStore{}, list)
}

func contextWithMailerSourceName(parent context.Context, mailerSourceName string) context.Context {
	return context.WithValue(parent, ctxMailerSourceName{}, mailerSourceName)
}

func contextWithNotificationEmail(parent context.Context, email string) context.Context {
	return context.WithValue(parent, ctxNotificationEmail{}, email)
}

func contextWithDataDir(parent context.Context, dirname string) context.Context {
	return context.WithValue(parent, ctxDataDir{}, dirname)
}

func contextWithVisitsInfoStore(parent context.Context, list *bome.JSONList) context.Context {
	return context.WithValue(parent, ctxVisitsInfo{}, list)
}

func contextWithResFS(parent context.Context, fs http.FileSystem) context.Context {
	return context.WithValue(parent, ctxResFS{}, fs)
}

func contextWithExternalResDir(parent context.Context, dirname string) context.Context {
	return context.WithValue(parent, ctxResDir{}, dirname)
}

func messages(ctx context.Context) *bome.JSONList {
	o := ctx.Value(ctxMessagesStore{})
	if o == nil {
		return nil
	}
	return o.(*bome.JSONList)
}

func visitsInfoStore(ctx context.Context) *bome.JSONList {
	o := ctx.Value(ctxVisitsInfo{})
	if o == nil {
		return nil
	}
	return o.(*bome.JSONList)
}

func resFS(ctx context.Context) http.FileSystem {
	o := ctx.Value(ctxResFS{})
	if o == nil {
		return nil
	}
	return o.(http.FileSystem)
}

func externalResDir(ctx context.Context) string {
	o := ctx.Value(ctxResDir{})
	if o == nil {
		return ""
	}
	return o.(string)
}

func notificationEmail(ctx context.Context) string {
	o := ctx.Value(ctxNotificationEmail{})
	if o == nil {
		return ""
	}
	return o.(string)
}

func mailerSource(ctx context.Context) string {
	o := ctx.Value(ctxMailerSourceName{})
	if o == nil {
		return ""
	}
	return o.(string)
}

func appDataDir(ctx context.Context) string {
	o := ctx.Value(ctxDataDir{})
	if o == nil {
		return ""
	}
	return o.(string)
}
