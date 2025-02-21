package main

import (
	"fmt"
	"log"

	"github.com/computerdane/gears"
)

func fetchForgePromotions() ForgePromotions {
	loc := fmt.Sprintf(
		"%s/maven/net/minecraftforge/forge/promotions_slim.json",
		gears.StringValue("forge-files-url"),
	)
	return fetchJson[ForgePromotions](loc)
}

func getForgeVersionString(mcVersion string) string {
	version := gears.StringValue("forge-version")
	versionString := mcVersion + "-" + version
	if version == "recommended" || version == "latest" {
		promos := fetchForgePromotions()
		value, exists := promos.Promos[versionString]
		if !exists {
			log.Fatal("Could not find Forge version: ", versionString)
		}
		versionString = mcVersion + "-" + value
	}
	return versionString
}

func getForgeDownloadUrl(mcVersion string) string {
	versionString := getForgeVersionString(mcVersion)
	return fmt.Sprintf(
		"%s/net/minecraftforge/forge/%s/forge-%s-installer.jar",
		gears.StringValue("forge-maven-url"),
		versionString,
		versionString,
	)
}
