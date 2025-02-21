package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/computerdane/dials"
)

func init() {
	dials.Add(&dials.Dial{
		Name:         "version",
		ValueType:    "string",
		DefaultValue: "latest",
		Shorthand:    "v",
	})
	dials.Add(&dials.Dial{
		Name:         "loader",
		ValueType:    "string",
		DefaultValue: "vanilla",
		Shorthand:    "l",
	})
	dials.Add(&dials.Dial{
		Name:         "forge-version",
		ValueType:    "string",
		DefaultValue: "recommended",
	})
	dials.Add(&dials.Dial{
		Name:      "overwrite",
		ValueType: "bool",
		Shorthand: "O",
	})
	dials.Add(&dials.Dial{
		Name:         "modrinth-modpack",
		Shorthand:    "M",
		ValueType:    "string",
		DefaultValue: "",
	})
	dials.Add(&dials.Dial{
		Name:         "modrinth-mod",
		Shorthand:    "m",
		ValueType:    "strings",
		DefaultValue: []string{},
	})
	dials.Add(&dials.Dial{
		Name:         "version-manifest-url",
		ValueType:    "string",
		DefaultValue: "https://launchermeta.mojang.com/mc/game/version_manifest.json",
	})
	dials.Add(&dials.Dial{
		Name:         "modrinth-api-url",
		ValueType:    "string",
		DefaultValue: "https://api.modrinth.com/v2",
	})
	dials.Add(&dials.Dial{
		Name:         "forge-files-url",
		ValueType:    "string",
		DefaultValue: "https://files.minecraftforge.net",
	})
	dials.Add(&dials.Dial{
		Name:         "forge-maven-url",
		ValueType:    "string",
		DefaultValue: "https://maven.minecraftforge.net",
	})

	dials.AddHomeConfigFile(".config/mc-quick/config.json")

	if configFile, exists := os.LookupEnv("CONFIG_FILE"); exists {
		fmt.Println(configFile)
		dials.AddConfigFile(configFile)
	}
}

func installModrinthModpack(mcVersion string) {
	printStep("Installing Modrinth modpack")

	modpack := dials.StringValue("modrinth-modpack")
	if modpack == "" {
		return
	}

	modpackFile := fetchModrinthProjectPrimaryFile(modpack, mcVersion)
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

func installModrinthMods(mcVersion string) {
	printStep("Installing Modrinth mods")

	for _, mod := range dials.StringValues("modrinth-mod") {
		file := fetchModrinthProjectPrimaryFile(mod, mcVersion)
		download("mods/"+file.Filename, file.Url, file.Hashes.Sha1)
	}
}

func install() {
	printStep("Installing server")

	mcVersionManifest := fetchJson[McVersionManifest](dials.StringValue("version-manifest-url"))

	mcVersion := dials.StringValue("version")
	switch mcVersion {
	case "latest":
		mcVersion = mcVersionManifest.Latest.Release
	case "latest-snapshot":
		mcVersion = mcVersionManifest.Latest.Snapshot
	}

	mcVersionMetaUrl := ""
	for _, entry := range mcVersionManifest.Versions {
		if entry.Id == mcVersion {
			mcVersionMetaUrl = entry.Url
		}
	}
	if mcVersionMetaUrl == "" {
		log.Fatalf("Could not find Minecraft version %s", mcVersion)
	}

	switch dials.StringValue("loader") {
	case "vanilla":
		mcVersionMeta := fetchJson[McVersionMetadata](mcVersionMetaUrl)
		download("server.jar", mcVersionMeta.Downloads.Server.Url, mcVersionMeta.Downloads.Server.Sha1)
	case "fabric":
		run("fabric-installer", "server", "-downloadMinecraft", "-mcversion", mcVersion)

		installModrinthModpack(mcVersion)
		installModrinthMods(mcVersion)
	case "forge":
		loc := getForgeDownloadUrl(mcVersion)
		download("forge-installer.jar", loc, "")

		run("java", "-jar", "forge-installer.jar", "--installServer")

		installModrinthModpack(mcVersion)
		installModrinthMods(mcVersion)
	}
}

func start() {
	printStep("Starting server")

	switch dials.StringValue("loader") {
	case "vanilla":
		run("java", "-jar", "server.jar", "--nogui")
	case "fabric":
		run("java", "-jar", "fabric-server-launch.jar", "--nogui")
	case "forge":
		run("./run.sh", "--nogui")
	}
}

func main() {
	dials.Load()

	args := dials.Positionals()
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

	loader := dials.StringValue("loader")
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
