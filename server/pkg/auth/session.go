package auth

import (
	"TaskHub/global"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func InitSessionStore() {
	store = sessions.NewCookieStore([]byte("recomend-32bytes-at-least"))
	store.Options = &sessions.Options{
		Path:     global.Cfg.Session.Path,
		MaxAge:   global.Cfg.Session.MaxAge,
		Secure:   global.Cfg.Session.Secure,
		HttpOnly: global.Cfg.Session.HttpOnly,
		Domain:   global.Cfg.Session.Domain,
		SameSite: global.Cfg.Session.SameSite,
	}
}

func SetSession(ctx *gin.Context, key string, value interface{}) error {

	session, err := store.Get(ctx.Request, "TaskHub")
	if err != nil {
		return err
	}

	session.Values[key] = value

	return session.Save(ctx.Request, ctx.Writer)
}

func GetSession(ctx *gin.Context, key string) interface{} {

	session, err := store.Get(ctx.Request, "TaskHub")
	if err != nil {
		return err
	}

	return session.Values[key]
}

func DelSession(ctx *gin.Context, key string) error {

	session, err := store.Get(ctx.Request, "TaskHub")
	if err != nil {
		return err
	}

	delete(session.Values, key)
	return session.Save(ctx.Request, ctx.Writer)
}
