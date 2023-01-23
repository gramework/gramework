package main

import (
	"os"

	"github.com/gramework/gramework"
	"github.com/gramework/gramework/behind/akamai"
)

func main() {
	app := gramework.New()

	csv, err := os.ReadFile("./testdata/all-cidr-blocks.csv")
	must(err)
	akamCIDR, err := akamai.ParseCIDRBlocksCSV(csv, true, true)
	must(err)
	app.Behind(akamai.New(akamai.CIDRBlocks(akamCIDR)))

	app.GET("/", func(ctx *gramework.Context) {
		ctx.WriteString(ctx.RemoteIP().String())
	})

	app.ListenAndServe()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
