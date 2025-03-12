package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/computerdane/gears"
	"github.com/computerdane/mc-quick/lib"
)

var Version string

//go:embed USAGE
var usage string

func printUsage() {
	fmt.Print(usage)
	gears.PrintUsage()
	os.Exit(1)
}

func init() {
	gears.Add(&gears.Flag{
		Name:        "version",
		ValueType:   "bool",
		Shorthand:   "v",
		Description: "Output the current version.",
	})
	gears.Add(&gears.Flag{
		Name:        "help",
		ValueType:   "bool",
		Shorthand:   "h",
		Description: "Show this help menu.",
	})
	gears.Add(&gears.Flag{
		Name:        "fish-completions",
		ValueType:   "bool",
		Description: "Generate fish completions for this command",
	})
	gears.Add(&gears.Flag{
		Name:         "mc-version",
		ValueType:    "string",
		DefaultValue: "latest",
		Shorthand:    "V",
		Description:  "The Minecraft version to use.",
	})
	gears.Add(&gears.Flag{
		Name:         "loader",
		ValueType:    "string",
		DefaultValue: "vanilla",
		Shorthand:    "L",
		Description:  "The mod loader to use; vanilla if no mods desired. Accepted values: vanilla fabric forge.",
	})
	gears.Add(&gears.Flag{
		Name:         "forge-version",
		ValueType:    "string",
		DefaultValue: "recommended",
		Description:  "The Forge version to use.",
	})
	gears.Add(&gears.Flag{
		Name:        "overwrite",
		ValueType:   "bool",
		Shorthand:   "O",
		Description: "If true, will re-download all files regardless of if their checksums are already valid.",
	})
	gears.Add(&gears.Flag{
		Name:         "modrinth-modpack",
		Shorthand:    "M",
		ValueType:    "string",
		DefaultValue: "",
		Description:  "Specify a Modrinth modpack to install.",
	})
	gears.Add(&gears.Flag{
		Name:         "modrinth-mod",
		Shorthand:    "m",
		ValueType:    "strings",
		DefaultValue: []string{},
		Description:  "Specify a Modrinth mod to install. Can be used multiple times to add multiple mods.",
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
	gears.Add(&gears.Flag{
		Name:         "fabric-installer-url",
		ValueType:    "string",
		DefaultValue: "https://maven.fabricmc.net/net/fabricmc/fabric-installer/1.0.1/fabric-installer-1.0.1.jar",
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

	modpackFile := lib.FetchModrinthProjectPrimaryFile(modpack)
	if !strings.HasSuffix(modpackFile.Filename, ".mrpack") {
		log.Fatalf("Modpack %s is not a .mrpack", modpack)
	}

	download(modpackFile.Filename, modpackFile.Url, modpackFile.Hashes.Sha1)
	run("unzip", "-qo", modpackFile.Filename)
	run("chmod", "-R", "700", "modrinth.index.json", "overrides")

	index := lib.JsonFile[lib.ModrinthModpackIndex]("modrinth.index.json")
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
		file := lib.FetchModrinthProjectPrimaryFile(mod)
		download("mods/"+file.Filename, file.Url, file.Hashes.Sha1)
	}
}

func install() {
	printStep("Installing server")

	mcVersionManifest := lib.FetchJson[lib.McVersionManifest](gears.StringValue("mc-version-manifest-url"))

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
		mcVersionMeta := lib.FetchJson[lib.McVersionMetadata](mcVersionMetaUrl)
		download("server.jar", mcVersionMeta.Downloads.Server.Url, mcVersionMeta.Downloads.Server.Sha1)
	case "fabric":
		download("fabric-installer.jar", gears.StringValue("fabric-installer-url"), "")
		run("java", "-jar", "fabric-installer.jar", "server", "-downloadMinecraft", "-mcversion", mcVersion)

		installModrinthModpack()
		installModrinthMods()
	case "forge":
		loc := lib.GetForgeDownloadUrl()
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

	if gears.BoolValue("help") {
		printUsage()
	}

	if gears.BoolValue("fish-completions") {
		fmt.Println(gears.FishCompletions("mc-quick"))
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
