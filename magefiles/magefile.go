//go:build mage
// +build mage

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
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

func (Version) New(ctx context.Context, inc string) error {
	mg.CtxDeps(ctx, Install)

	args := []string{"generate", "--auto-tag", "--auto-commit", "--auto-push"}

	currentVersion, _ := currentVersion(ctx)
	if currentVersion != nil {
		if inc == "patch" {
			args = append(args, "--version", fmt.Sprintf("v%s", currentVersion.IncPatch().String()))
		}
		if inc == "minor" {
			args = append(args, "--version", fmt.Sprintf("v%s", currentVersion.IncMinor().String()))
		}
		if inc == "major" {
			args = append(args, "--version", fmt.Sprintf("v%s", currentVersion.IncMajor().String()))
		}
	} else {
		args = append(args, "--version", "v0.0.1")
	}

	cmd := exec.Command("changeloguru", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
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
