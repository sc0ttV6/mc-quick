package lib

import (
	"fmt"
	"log"
	"net/url"

	"github.com/computerdane/gears"
)

func modrinthUrlSingleton(s string) string {
	return url.QueryEscape(fmt.Sprintf(`["%s"]`, s))
}

func FetchModrinthProjectVersions(slug string) ModrinthProjectVersions {
	loc := fmt.Sprintf(
		`%s/project/%s/version?loaders=%s&game_versions=%s`,
		gears.StringValue("modrinth-api-url"),
		slug,
		modrinthUrlSingleton(gears.StringValue("loader")),
		modrinthUrlSingleton(gears.StringValue("mc-version")),
	)
	return FetchJson[ModrinthProjectVersions](loc)
}

func FetchModrinthProjectPrimaryFile(slug string, versionNumber string) ModrinthProjectVersionsFile {
	versions := FetchModrinthProjectVersions(slug)
	if len(versions) == 0 {
		log.Fatalf("Could not find a valid version on Modrinth for %s for Minecraft %s", slug, gears.StringValue("mc-version"))
	}

	var version struct {
		VersionNumber string `json:"version_number"`
		Files         []ModrinthProjectVersionsFile
	}

	if versionNumber != "" {
		found := false
		for _, v := range versions {
			if v.VersionNumber == versionNumber {
				version = v
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("Could not find version %s on Modrinth for %s", versionNumber, slug)
		}
	} else {
		version = versions[0]
	}

	for _, file := range version.Files {
		if file.Primary {
			return file
		}
	}

	log.Fatalf("Found version on Modrinth for %s, but could not find a primary file", slug)
	return ModrinthProjectVersionsFile{}
}
