package lib

type McVersionManifest struct {
	Latest struct {
		Release  string
		Snapshot string
	}
	Versions []struct {
		Id  string
		Url string
	}
}

type McVersionMetadata struct {
	Downloads struct {
		Server struct {
			Sha1 string
			Url  string
		}
	}
}

type ModrinthProjectVersionsFile struct {
	Hashes struct {
		Sha1 string
	}
	Url      string
	Filename string
	Primary  bool
}

type ModrinthProjectVersions []struct {
	Files []ModrinthProjectVersionsFile
}

type ModrinthModpackIndexFile struct {
	Path   string
	Hashes struct {
		Sha1 string
	}
	Env struct {
		Client string
		Server string
	}
	Downloads []string
}

type ModrinthModpackIndex struct {
	Files []ModrinthModpackIndexFile
}

type ForgePromotions struct {
	Promos map[string]string
}
