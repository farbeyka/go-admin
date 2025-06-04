package beego

import (
	"net/http"
	"os"

	// add beego adapter
	_ "github.com/farbeyka/go-admin/adapter/beego2"
	// add mysql driver
	_ "github.com/farbeyka/go-admin/modules/db/drivers/mysql"
	// add postgresql driver
	_ "github.com/farbeyka/go-admin/modules/db/drivers/postgres"
	// add sqlite driver
	_ "github.com/farbeyka/go-admin/modules/db/drivers/sqlite"
	// add mssql driver
	_ "github.com/farbeyka/go-admin/modules/db/drivers/mssql"

	"github.com/farbeyka/go-admin/engine"
	"github.com/farbeyka/go-admin/modules/config"
	"github.com/farbeyka/go-admin/modules/language"
	"github.com/farbeyka/go-admin/plugins/admin"
	"github.com/farbeyka/go-admin/plugins/admin/modules/table"
	"github.com/farbeyka/go-admin/plugins/example"
	"github.com/farbeyka/go-admin/template"
	"github.com/farbeyka/go-admin/template/chartjs"
	"github.com/farbeyka/go-admin/tests/tables"
	"github.com/GoAdminGroup/themes/adminlte"
	"github.com/beego/beego/v2/server/web"
)

func internalHandler() http.Handler {

	app := web.NewHttpSever()

	eng := engine.Default()
	adminPlugin := admin.NewAdmin(tables.Generators)
	adminPlugin.AddGenerator("user", tables.GetUserTable)

	examplePlugin := example.NewExample()

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin, examplePlugin).Use(app); err != nil {
		panic(err)
	}

	template.AddComp(chartjs.NewChart())

	eng.HTML("GET", "/admin", tables.GetContent)

	app.Cfg.Listen.HTTPAddr = "127.0.0.1"
	app.Cfg.Listen.HTTPPort = 9087

	return app.Handlers
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {

	app := web.NewHttpSever()

	eng := engine.Default()
	adminPlugin := admin.NewAdmin(gens)

	if err := eng.AddConfig(&config.Config{
		Databases: dbs,
		UrlPrefix: "admin",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language:    language.EN,
		IndexUrl:    "/",
		Debug:       true,
		ColorScheme: adminlte.ColorschemeSkinBlack,
	}).
		AddPlugins(adminPlugin).Use(app); err != nil {
		panic(err)
	}

	template.AddComp(chartjs.NewChart())

	eng.HTML("GET", "/admin", tables.GetContent)

	app.Cfg.Listen.HTTPAddr = "127.0.0.1"
	app.Cfg.Listen.HTTPPort = 9087

	return app.Handlers
}
