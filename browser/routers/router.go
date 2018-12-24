package routers

import (
	"github.com/astaxie/beego"
	. "github.com/axengine/btchain/browser/controllers"
)

func init() {
	beego.Router("/", &WebController{}, "get:Index")
	webNs := beego.NewNamespace("/view",
		beego.NSRouter("/blocks/latest", &WebController{}, "get:Latest"),
		beego.NSRouter("/blocks/hash/:hash", &WebController{}, "get:Block"),

		beego.NSRouter("/txs/page", &WebController{}, "get:TxsPage"),
		beego.NSRouter("/trans/txid/:txid", &WebController{}, "get:Action"),
		beego.NSRouter("/trans/detail/:hash", &WebController{}, "get:TxByTxHash"),
		beego.NSRouter("/accounts/:address", &WebController{}, "get:TxOutByAddress"),
		beego.NSRouter("/accounts/:address/income", &WebController{}, "get:TxInByAddress"),
		beego.NSRouter("/accounts/:address/payout", &WebController{}, "get:TxOutByAddress"),

		beego.NSRouter("/search/:hash", &WebController{}, "get:Search"),
	)
	beego.AddNamespace(webNs)
	beego.SetStaticPath("/assets", "static")
}
