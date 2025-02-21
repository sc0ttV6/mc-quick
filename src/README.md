# mc-quick

Get a Minecraft server running ASAP

# Features

- Supports vanilla, Fabric, and Forge servers
- Supports mods and modpacks from Modrinth
- Supports configuration with JSON, environment variables, and command-line arguments
- Skips re-downloading existing files that have valid checksums

# Dependencies

Please ensure the following binaries are installed on your system:

- java
- rsync
- unzip

# Usage Examples

Start a vanilla server on the latest version:

```sh
  mc-quick install
  mc-quick start
```

Start a Fabric server with a few mods:

```sh
  mc-quick --loader fabric install -m fabric-api -m simple-voice-chat -m no-chat-reports
  mc-quick --loader fabric start
```

Start a Forge server with a modpack:

```sh
  mc-quick -L forge install -V 1.20.1 -M cave-horror-project-modpack
  mc-quick -L forge start
```

Configure using environment variables:

```sh
  export LOADER=fabric
  export MC_VERSION=1.19.4
  mc-quick install
  mc-quick start
```

Configure using JSON:

```json
  {
    "loader": "fabric",
    "mc-version": "1.19.4",
    "modrinth-mod": ["fabric-api", "simple-voice-chat", "no-chat-reports"]
  }
```

then,

```sh
  export CONFIG_FILE=/path/to/config.json
  mc-quick install
  mc-quick start
```
