package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/computerdane/gears"
)

var Version string

func init() {
	gears.Add(&gears.Flag{
		Name:      "version",
		ValueType: "bool",
		Shorthand: "v",
	})
	gears.Add(&gears.Flag{
		Name:         "mc-version",
		ValueType:    "string",
		DefaultValue: "latest",
		Shorthand:    "V",
	})
	gears.Add(&gears.Flag{
		Name:         "loader",
		ValueType:    "string",
		DefaultValue: "vanilla",
		Shorthand:    "L",
	})
	gears.Add(&gears.Flag{
		Name:         "forge-version",
		ValueType:    "string",
		DefaultValue: "recommended",
	})
	gears.Add(&gears.Flag{
		Name:      "overwrite",
		ValueType: "bool",
		Shorthand: "O",
	})
	gears.Add(&gears.Flag{
		Name:         "modrinth-modpack",
		Shorthand:    "M",
		ValueType:    "string",
		DefaultValue: "",
	})
	gears.Add(&gears.Flag{
		Name:         "modrinth-mod",
		Shorthand:    "m",
		ValueType:    "strings",
		DefaultValue: []string{},
	})
	gears.Add(&gears.Flag{
		Name:         "mc-version-manifest-url",
		ValueType:    "string",
		DefaultValue: "https://launchermeta.mojang.com/mc/game/version_manifest.json",
	})
	gears.Add(&gears.Flag{
		Name:         "modrinth-api-url",
		ValueType:    "string",
		DefaultValue: "https://api.modrinth.com/v2",
	})
	gears.Add(&gears.Flag{
		Name:         "forge-files-url",
		ValueType:    "string",
		DefaultValue: "https://files.minecraftforge.net",
	})
	gears.Add(&gears.Flag{
		Name:         "forge-maven-url",
		ValueType:    "string",
		DefaultValue: "https://maven.minecraftforge.net",
	})

	gears.AddHomeConfigFile(".config/mc-quick/config.json")

	if configFile, exists := os.LookupEnv("CONFIG_FILE"); exists {
		gears.AddConfigFile(configFile)
	}
}

func installModrinthModpack() {
	printStep("Installing Modrinth modpack")

	modpack := gears.StringValue("modrinth-modpack")
	if modpack == "" {
		return
	}

	modpackFile := fetchModrinthProjectPrimaryFile(modpack)
	if !strings.HasSuffix(modpackFile.Filename, ".mrpack") {
		log.Fatalf("Modpack %s is not a .mrpack", modpack)
	}

	download(modpackFile.Filename, modpackFile.Url, modpackFile.Hashes.Sha1)
	run("unzip", "-qo", modpackFile.Filename)
	run("chmod", "-R", "700", "modrinth.index.json", "overrides")

	index := jsonFile[ModrinthModpackIndex]("modrinth.index.json")
	for _, file := range index.Files {
		if file.Env.Server == "required" {
			if len(file.Downloads) == 0 {
				fmt.Printf("Warning: Modpack mod %s does not have a download URL", file.Path)
				continue
			}
			download(file.Path, file.Downloads[0], file.Hashes.Sha1)
		}
	}

	fmt.Println("Applying overrides...")
	run("rsync", "-aI", "overrides/", ".")
}

func installModrinthMods() {
	printStep("Installing Modrinth mods")

	for _, mod := range gears.StringValues("modrinth-mod") {
		file := fetchModrinthProjectPrimaryFile(mod)
		download("mods/"+file.Filename, file.Url, file.Hashes.Sha1)
	}
}

func install() {
	printStep("Installing server")

	mcVersionManifest := fetchJson[McVersionManifest](gears.StringValue("mc-version-manifest-url"))

	switch gears.StringValue("mc-version") {
	case "latest":
		gears.SetValue("mc-version", mcVersionManifest.Latest.Release)
	case "latest-snapshot":
		gears.SetValue("mc-version", mcVersionManifest.Latest.Snapshot)
	}
	mcVersion := gears.StringValue("mc-version")

	mcVersionMetaUrl := ""
	for _, entry := range mcVersionManifest.Versions {
		if entry.Id == mcVersion {
			mcVersionMetaUrl = entry.Url
		}
	}
	if mcVersionMetaUrl == "" {
		log.Fatalf("Could not find Minecraft version %s", mcVersion)
	}

	switch gears.StringValue("loader") {
	case "vanilla":
		mcVersionMeta := fetchJson[McVersionMetadata](mcVersionMetaUrl)
		download("server.jar", mcVersionMeta.Downloads.Server.Url, mcVersionMeta.Downloads.Server.Sha1)
	case "fabric":
		run("fabric-installer", "server", "-downloadMinecraft", "-mcversion", mcVersion)

		installModrinthModpack()
		installModrinthMods()
	case "forge":
		loc := getForgeDownloadUrl()
		download("forge-installer.jar", loc, "")

		run("java", "-jar", "forge-installer.jar", "--installServer")

		installModrinthModpack()
		installModrinthMods()
	}
}

func start() {
	printStep("Starting server")

	switch gears.StringValue("loader") {
	case "vanilla":
		run("java", "-jar", "server.jar", "--nogui")
	case "fabric":
		run("java", "-jar", "fabric-server-launch.jar", "--nogui")
	case "forge":
		run("./run.sh", "--nogui")
	}
}

func main() {
	gears.Load()

	if gears.BoolValue("version") {
		if Version == "" {
			Version = "unknown version"
		}
		fmt.Printf("mc-quick  %s\n", Version)
		os.Exit(0)
	}

	args := gears.Positionals()
	if len(args) == 0 {
		printUsage()
	} else if len(args) == 1 {
		if args[0] != "install" &&
			args[0] != "start" {
			fmt.Println("Invalid command: ", args[0])
			printUsage()
		}
	} else {
		fmt.Println("Error: Too many arguments!")
		printUsage()
	}

	loader := gears.StringValue("loader")
	if loader != "vanilla" &&
		loader != "fabric" &&
		loader != "forge" {
		log.Fatalf("Loader must be either vanilla, fabric, or forge")
	}

	switch args[0] {
	case "install":
		install()
	case "start":
		start()
	}
}
