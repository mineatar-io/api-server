package routes

import (
	"net/http"

	"github.com/mineatar-io/api-server/src/util"
	"github.com/valyala/fasthttp"
)

func UUIDHandler(ctx *fasthttp.RequestCtx) {
	user := ctx.UserValue("user").(string)

	uuid, ok, err := util.LookupUUID(user)

	if err != nil {
		util.WriteError(ctx, err, http.StatusInternalServerError)

		return
	}

	if !ok {
		util.WriteError(ctx, nil, http.StatusNotFound)

		return
	}

	ctx.SetBodyString(uuid)
}
