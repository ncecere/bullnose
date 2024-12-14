package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/ncecere/bullnose/internal/config"
	"github.com/ncecere/bullnose/internal/scraper"
)

// Version represents the current version of bullnose.
// It is set during build time using -ldflags.
var Version = "dev"

var cfgFile string

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:     "bullnose [flags] [urls...]",
	Short:   "A web scraper that converts web pages to clean markdown",
	Version: Version,
	RunE:    run,
	Args:    cobra.ArbitraryArgs,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./bullnose.yaml)")
	rootCmd.Flags().StringP("output", "o", "./scraped-content", "output directory for scraped content")
	rootCmd.Flags().IntP("depth", "d", 3, "maximum depth to follow links")
	rootCmd.Flags().IntP("parallel", "p", 8, "number of parallel scraping actions")
	rootCmd.Flags().BoolP("restrict-domain", "r", true, "only follow links within starting domain")
	rootCmd.Flags().String("rescrape-after", "12h", "only rescrape after this duration (format: Xm, Xh, Xd)")
	rootCmd.Flags().BoolP("force", "f", false, "force rescrape regardless of time")
	rootCmd.Flags().Bool("debug", false, "enable debug logging")
	rootCmd.Flags().StringSlice("ignore", []string{}, "URLs or patterns to ignore")
}

func run(cmd *cobra.Command, args []string) error {
	// Load config from file and flags
	cfg, err := loadConfig(cmd, args)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Validate that we have at least one URL
	if len(cfg.URLs) == 0 && len(args) == 0 {
		return fmt.Errorf("at least one URL must be provided")
	}

	// Add command line URLs to config URLs
	if len(args) > 0 {
		cfg.URLs = append(cfg.URLs, args...)
	}

	// Create and run scraper
	s, err := scraper.New(cfg)
	if err != nil {
		return fmt.Errorf("error creating scraper: %w", err)
	}

	if err := s.Start(); err != nil {
		return fmt.Errorf("error running scraper: %w", err)
	}

	return nil
}

func loadConfig(cmd *cobra.Command, args []string) (*config.Config, error) {
	// Load config from file first
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return nil, err
	}

	// Override with command line flags
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return nil, fmt.Errorf("error getting output flag: %w", err)
	}
	if cmd.Flags().Changed("output") {
		cfg.Output = output
	}

	depth, err := cmd.Flags().GetInt("depth")
	if err != nil {
		return nil, fmt.Errorf("error getting depth flag: %w", err)
	}
	if cmd.Flags().Changed("depth") {
		cfg.Depth = depth
	}

	parallel, err := cmd.Flags().GetInt("parallel")
	if err != nil {
		return nil, fmt.Errorf("error getting parallel flag: %w", err)
	}
	if cmd.Flags().Changed("parallel") {
		cfg.Parallel = parallel
	}

	restrictDomain, err := cmd.Flags().GetBool("restrict-domain")
	if err != nil {
		return nil, fmt.Errorf("error getting restrict-domain flag: %w", err)
	}
	if cmd.Flags().Changed("restrict-domain") {
		cfg.RestrictDomain = restrictDomain
	}

	if cmd.Flags().Changed("rescrape-after") {
		rescrapeAfter, err := cmd.Flags().GetString("rescrape-after")
		if err != nil {
			return nil, fmt.Errorf("error getting rescrape-after flag: %w", err)
		}
		duration, err := time.ParseDuration(rescrapeAfter)
		if err != nil {
			return nil, fmt.Errorf("invalid rescrape-after duration: %w", err)
		}
		cfg.RescrapeAfter = duration
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return nil, fmt.Errorf("error getting force flag: %w", err)
	}
	if cmd.Flags().Changed("force") {
		cfg.Force = force
	}

	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return nil, fmt.Errorf("error getting debug flag: %w", err)
	}
	if cmd.Flags().Changed("debug") {
		cfg.Debug = debug
	}

	ignore, err := cmd.Flags().GetStringSlice("ignore")
	if err != nil {
		return nil, fmt.Errorf("error getting ignore flag: %w", err)
	}
	if cmd.Flags().Changed("ignore") {
		cfg.Ignore = ignore
	}

	return cfg, nil
}
