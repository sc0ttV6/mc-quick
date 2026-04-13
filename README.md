# mc-quick

A NixOS module for quickly setting up Minecraft servers with support for vanilla, Fabric, and Forge, including Modrinth modpacks.

Built on [mc-quick](https://github.com/computerdane/mc-quick) (v1.0.5).

## Usage

Add the flake input:

```nix
# flake.nix
{
  inputs = {
    mc-quick = {
      type = "sourcehut";
      owner = "~scotty";
      repo = "mc-quick";
    };
  };
}
```

Import the module in your NixOS configuration:

```nix
{
  imports = [
    mc-quick.nixosModules.default
  ];
}
```

## Examples

### Vanilla server

```nix
services.mc-quick.vanilla = {
  enable = true;
  autoStart = true;
  acceptEula = true;
  openFirewall = true;
  serverProperties = {
    motd = "A vanilla Minecraft server";
  };
};
```

### Modded server (Pixelmon)

```nix
services.mc-quick.pixelmon = {
  enable = true;
  autoStart = true;
  acceptEula = true;
  openFirewall = true;
  mcVersion = "1.21.1";
  loader = "forge";
  modrinthModpack = "the-pixelmon-modpack";
  serverProperties = {
    motd = "Pixelmon server";
  };
};
```

### Fabric server with Modrinth mods

```nix
services.mc-quick.fabric = {
  enable = true;
  autoStart = true;
  acceptEula = true;
  openFirewall = true;
  loader = "fabric";
  modrinthMods = [
    "lithium"
    "fabric-api"
  ];
  serverProperties = {
    motd = "Fabric server";
  };
};
```

## Options

| Option | Default | Description |
|---|---|---|
| `enable` | `false` | Enable the Minecraft server |
| `autoStart` | `false` | Start the server on boot |
| `mcVersion` | `"1.21.4"` | Minecraft version |
| `loader` | `"vanilla"` | Mod loader: `vanilla`, `fabric`, or `forge` |
| `forgeVersion` | `"recommended"` | Forge version |
| `acceptEula` | `false` | Accept the Minecraft EULA |
| `port` | `25565` | Server port |
| `rconPort` | `25575` | RCON port |
| `openFirewall` | `false` | Open firewall for the server port |
| `modrinthModpack` | `""` | Modrinth modpack project slug |
| `modrinthMods` | `[]` | List of Modrinth mod project slugs |
| `serverProperties` | `{ motd = "Powered by NixOS!"; }` | server.properties values |
| `ops` | `[]` | Server operators (uuid, name, level, bypassesPlayerLimit) |
| `whitelist` | `[]` | Whitelisted players (uuid, name) |
| `enableWhitelist` | `false` | Enable whitelisting |
| `javaPackage` | `temurin-jre-bin` | Java package to use |
| `files` | `[]` | Extra files to write in the server directory (path, text) |
