package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
	_ "github.com/nakagami/firebirdsql"
)

type Template struct {
	templates *template.Template
}

type TemplateRenderer struct {
	templates *template.Template
}

type SqlTableContent struct {
	ID_AGR int
	PODR   int
	UID    int
	OIK    int
	OTD    int
	AGR    int
	UIK    int
	OTM    int
	DTA    string
	SPZ    int
	KMP    int
	DS     string
	DRED   string
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	db := dbConn()
	query := `select * from agree where uid>517600`

	//	selDB, err := db.Query("SELECT SALE_ID, DELIVERY_DATE, NAME FROM PREORDER") This works
	selDB, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	log.Println(selDB)
	var queryContent SqlTableContent
	var SqlTableContentArray []SqlTableContent
	for selDB.Next() {

		var id_agr int
		var podr int
		var uid int
		var oik int
		var otd int
		var agr int
		var uik int
		var otm int
		var dta string
		var spz int
		var kmp int
		var ds string
		var dred string

		err = selDB.Scan(&id_agr, &podr, &uid, &oik, &otd, &agr, &uik, &otm, &dta, &spz, &kmp, &ds, &dred)
		if err != nil {
			panic(err.Error())
		}
		queryContent.ID_AGR = id_agr
		queryContent.PODR = podr
		queryContent.UID = uid
		queryContent.OIK = oik
		queryContent.OTD = otd
		queryContent.AGR = agr
		queryContent.UIK = uik
		queryContent.OTM = otm
		queryContent.DTA = dta
		queryContent.SPZ = spz
		queryContent.KMP = kmp
		queryContent.DS = ds
		queryContent.DRED = dred

		SqlTableContentArray = append(SqlTableContentArray, queryContent)

	}
	// tmpl.ExecuteTemplate(w, "Show", SqlTableContentArray)
	log.Println(SqlTableContentArray)
	defer db.Close()

	e := echo.New()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("form/*.html")),
	}
	e.Renderer = renderer

	// Named route "foobar"
	// e.GET("/show", func(c echo.Context) error {
	// 	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
	// 		"name": "Dolly!",
	// 		"data": SqlTableContentArray,
	// 	})
	// }).Name = "foobar"

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", SqlTableContentArray)
	}).Name = "foobar"

	e.Logger.Fatal(e.Start(":8000"))
}

func dbConn() (db *sql.DB) {
	db, err := sql.Open("firebirdsql", "sysdba:masterkey@192.168.0.15/c:/PublicFolders/Services/Arena/DB/ARENA.GDB")
	if err != nil {
		panic(err.Error())
	}
	return db
}

// var tmpl = template.Must(template.ParseGlob("form/*"))
