package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/ncecere/bullnose/internal/config"
	"github.com/ncecere/bullnose/internal/scraper"
)

var (
	Version = "dev"
	cfgFile string
)

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
	if cmd.Flags().Changed("output") {
		output, _ := cmd.Flags().GetString("output")
		cfg.Output = output
	}
	if cmd.Flags().Changed("depth") {
		depth, _ := cmd.Flags().GetInt("depth")
		cfg.Depth = depth
	}
	if cmd.Flags().Changed("parallel") {
		parallel, _ := cmd.Flags().GetInt("parallel")
		cfg.Parallel = parallel
	}
	if cmd.Flags().Changed("restrict-domain") {
		restrictDomain, _ := cmd.Flags().GetBool("restrict-domain")
		cfg.RestrictDomain = restrictDomain
	}
	if cmd.Flags().Changed("rescrape-after") {
		rescrapeAfter, _ := cmd.Flags().GetString("rescrape-after")
		duration, err := time.ParseDuration(rescrapeAfter)
		if err != nil {
			return nil, fmt.Errorf("invalid rescrape-after duration: %w", err)
		}
		cfg.RescrapeAfter = duration
	}
	if cmd.Flags().Changed("force") {
		force, _ := cmd.Flags().GetBool("force")
		cfg.Force = force
	}
	if cmd.Flags().Changed("debug") {
		debug, _ := cmd.Flags().GetBool("debug")
		cfg.Debug = debug
	}
	if cmd.Flags().Changed("ignore") {
		ignore, _ := cmd.Flags().GetStringSlice("ignore")
		cfg.Ignore = ignore
	}

	return cfg, nil
}
