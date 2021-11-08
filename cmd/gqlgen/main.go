// +build ignore
package main

import (
	"log"

	"github.com/hayashiki/audiy-api/src/graph/plugin/modelgen"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"

	"github.com/99designs/gqlgen/plugin/resolvergen"
)

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		log.Fatalln("cannot find gqlgen.yml anywhere from current working directory:", err)
	}

	log.Println("Regenerating GraphQL models and resolvers...")
	err = api.Generate(cfg,
		api.NoPlugins(),
		api.AddPlugin(modelgen.New()),
		api.AddPlugin(resolvergen.New()),
	)

	if err != nil {
		log.Fatalln("gqlgen generation error:", err)
	}
}
