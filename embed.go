package swamp

import "embed"

//go:embed templates/**
//go:embed static/**
//go:embed swamp_repos.yml
var appFS embed.FS
