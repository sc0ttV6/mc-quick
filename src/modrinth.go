package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/computerdane/gears"
)

func modrinthUrlSingleton(s string) string {
	return url.QueryEscape(fmt.Sprintf(`["%s"]`, s))
}

func fetchModrinthProjectVersions(slug string, mcVersion string) ModrinthProjectVersions {
	loc := fmt.Sprintf(
		`%s/project/%s/version?loaders=%s&game_versions=%s`,
		gears.StringValue("modrinth-api-url"),
		slug,
		modrinthUrlSingleton(gears.StringValue("loader")),
		modrinthUrlSingleton(mcVersion),
	)
	return fetchJson[ModrinthProjectVersions](loc)
}

func fetchModrinthProjectPrimaryFile(slug string, mcVersion string) ModrinthProjectVersionsFile {
	versions := fetchModrinthProjectVersions(slug, mcVersion)
	if len(versions) == 0 {
		log.Fatalf("Could not find a valid version on Modrinth for %s for Minecraft %s", slug, mcVersion)
	}

	version := versions[0]
	for _, file := range version.Files {
		if file.Primary {
			return file
		}
	}

	log.Fatalf("Found version on Modrinth for %s, but could not find a primary file", slug)
	return ModrinthProjectVersionsFile{}
}
