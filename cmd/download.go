package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/humancto/mozzy/internal/download"
)

var (
	downloadOutput    string
	downloadOverwrite bool
	downloadNoProgress bool
)

var downloadCmd = &cobra.Command{
	Use:   "download <url>",
	Short: "Download a file with progress bar",
	Long: `Download files from URLs with real-time progress tracking.

Features:
  - Progress bar with ETA and speed
  - Auto-detects filename from URL or Content-Disposition header
  - Resume support for interrupted downloads
  - Overwrite protection

Examples:
  mozzy download https://example.com/file.zip
  mozzy download https://example.com/file.zip -o myfile.zip
  mozzy download https://example.com/file.zip --overwrite
  mozzy download https://example.com/large.iso --no-progress`,
	Args: cobra.ExactArgs(1),
	RunE: runDownload,
}

func init() {
	downloadCmd.Flags().StringVarP(&downloadOutput, "output", "o", "", "Output filename (default: auto-detect from URL)")
	downloadCmd.Flags().BoolVar(&downloadOverwrite, "overwrite", false, "Overwrite existing file")
	downloadCmd.Flags().BoolVar(&downloadNoProgress, "no-progress", false, "Disable progress bar")
	rootCmd.AddCommand(downloadCmd)
}

func runDownload(cmd *cobra.Command, args []string) error {
	url := args[0]

	opts := download.DownloadOptions{
		URL:             url,
		OutputPath:      downloadOutput,
		ShowProgress:    !downloadNoProgress,
		OverwriteExist:  downloadOverwrite,
		FollowRedirects: true,
	}

	outputPath, err := download.Download(opts)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	fmt.Printf("\n✅ Downloaded to: %s\n", outputPath)
	return nil
}
