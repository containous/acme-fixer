package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/containous/acme-fixer/check"
	"github.com/containous/acme-fixer/traefikv1"
	"github.com/containous/acme-fixer/traefikv2"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// Config Command line configuration.
type Config struct {
	Filename   string
	DomainList string
	Version    bool
	DryRun     bool
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := Config{}

	rootCmd := &cobra.Command{
		Use:   "acme-fixer",
		Short: "Check and clean acme.json file.",
		Long:  `Check and clean acme.json file.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cfg.DomainList != "" {
				return nil
			}

			if cfg.Filename == "" {
				return errors.New("input is required")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.DomainList != "" {
				return checkDomainList(cfg.DomainList)
			}

			if cfg.Version {
				return traefikv2.Process(cfg.Filename, cfg.DryRun)
			}

			return traefikv1.Process(cfg.Filename, cfg.DryRun)
		},
	}

	flags := rootCmd.Flags()
	flags.StringVarP(&cfg.Filename, "input", "i", "", "The path of the acme.json.")
	flags.StringVarP(&cfg.DomainList, "domain-list", "l", "", "The path to a filename (list domains).")
	flags.BoolVarP(&cfg.DryRun, "dry-run", "d", false, "Dry run mode.")
	flags.BoolVar(&cfg.Version, "v2", false, "Use Traefik v2 format.")
	_ = flags.MarkHidden("domain-list")

	docCmd := &cobra.Command{
		Use:    "doc",
		Short:  "Generate documentation",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doc.GenMarkdownTree(rootCmd, "./docs")
		},
	}

	rootCmd.AddCommand(docCmd)

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display version",
		Run: func(_ *cobra.Command, _ []string) {
			displayVersion()
		},
	}

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func checkDomainList(filename string) error {
	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return err
	}

	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		text := fileScanner.Text()

		safe, err := check.IsSafe(strings.TrimSpace(text), false)
		if err != nil {
			return err
		}

		if !safe {
			fmt.Println(strings.TrimSpace(text))
		}
	}

	return nil
}
