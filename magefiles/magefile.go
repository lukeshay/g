//go:build mage
// +build mage

package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/Masterminds/semver/v3"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Version mg.Namespace

func (Version) Show(ctx context.Context) {
	mg.CtxDeps(ctx, Install)

	version, err := currentVersion(ctx)
	if err != nil {
		fmt.Printf("No version yet\n")
	} else {
		fmt.Printf("Version: v%s\n", version.String())
	}
}

func (Version) New(ctx context.Context) error {
	mg.CtxDeps(ctx, Install)

	return sh.RunV("changeloguru", "generate", "-i", "--auto-tag")
}

func currentVersion(ctx context.Context) (*semver.Version, error) {
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, "git", "describe", "--tags", "--abbrev=0", "--match=v[0-9]*", "HEAD")

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	out.Truncate(out.Len() - 1)

	return semver.NewVersion(out.String())
}

// Install CLIs
func Install() error {
	fmt.Println("Installing Deps...")

	return sh.RunV("go", "install", "github.com/haunt98/changeloguru/cmd/changeloguru@latest")
}
