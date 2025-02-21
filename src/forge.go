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

func getForgeVersionString() string {
	mcVersion := gears.StringValue("mc-version")
	forgeVersion := gears.StringValue("forge-version")
	versionString := mcVersion + "-" + forgeVersion
	if forgeVersion == "recommended" || forgeVersion == "latest" {
		promos := fetchForgePromotions()
		value, exists := promos.Promos[versionString]
		if !exists {
			log.Fatal("Could not find Forge version: ", versionString)
		}
		versionString = mcVersion + "-" + value
	}
	return versionString
}

func getForgeDownloadUrl() string {
	versionString := getForgeVersionString()
	return fmt.Sprintf(
		"%s/net/minecraftforge/forge/%s/forge-%s-installer.jar",
		gears.StringValue("forge-maven-url"),
		versionString,
		versionString,
	)
}
