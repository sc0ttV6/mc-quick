package lib

import (
	"fmt"
	"log"
	"strings"

	"github.com/computerdane/gears"
)

type NeoForgeVersions struct {
	Versions []string `json:"versions"`
}

func FetchNeoForgeVersions() NeoForgeVersions {
	loc := fmt.Sprintf(
		"%s/api/maven/versions/releases/net/neoforged/neoforge",
		gears.StringValue("neoforge-maven-url"),
	)
	return FetchJson[NeoForgeVersions](loc)
}

func GetNeoForgeVersion() string {
	neoforgeVersion := gears.StringValue("neoforge-version")
	if neoforgeVersion != "latest" {
		return neoforgeVersion
	}

	mcVersion := gears.StringValue("mc-version")
	parts := strings.SplitN(mcVersion, ".", 3)
	if len(parts) < 3 {
		log.Fatalf("Cannot derive NeoForge version prefix from Minecraft version %s", mcVersion)
	}
	prefix := parts[1] + "." + parts[2] + "."

	versions := FetchNeoForgeVersions()
	latest := ""
	for _, v := range versions.Versions {
		if strings.HasPrefix(v, prefix) {
			latest = v
		}
	}
	if latest == "" {
		log.Fatalf("Could not find a NeoForge version for Minecraft %s", mcVersion)
	}
	return latest
}

func GetNeoForgeDownloadUrl() string {
	version := GetNeoForgeVersion()
	return fmt.Sprintf(
		"%s/releases/net/neoforged/neoforge/%s/neoforge-%s-installer.jar",
		gears.StringValue("neoforge-maven-url"),
		version,
		version,
	)
}
