package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/humancto/mozzy/internal/ui"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update mozzy to the latest version",
	Long: `Check for updates and install the latest version of mozzy.

This command will:
- Check GitHub for the latest release
- Compare with your current version
- Download and install the update if available

Example:
  mozzy update`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get current version
		currentVersion := version
		fmt.Printf("Current version: %s\n\n", ui.InfoStyle.Render("v"+currentVersion))

		// Fetch latest version from GitHub
		fmt.Println("Checking for updates...")
		latestVersion, err := getLatestVersion()
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}

		fmt.Printf("Latest version:  %s\n\n", ui.SuccessStyle.Render("v"+latestVersion))

		// Compare versions
		if currentVersion == latestVersion {
			fmt.Println(ui.SuccessBanner("You're already on the latest version!"))
			return nil
		}

		// Prompt for update
		fmt.Println(ui.InfoBanner(fmt.Sprintf("New version available: v%s", latestVersion)))
		fmt.Println("\nUpdating mozzy...")

		// Run the install script
		installURL := "https://raw.githubusercontent.com/humancto/mozzy/main/install.sh"
		execCmd := exec.Command("bash", "-c", fmt.Sprintf("curl -fsSL %s | bash", installURL))
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		if err := execCmd.Run(); err != nil {
			return fmt.Errorf("update failed: %w", err)
		}

		fmt.Println("\n" + ui.SuccessBanner("Update completed successfully!"))
		fmt.Println(ui.InfoStyle.Render("💡 Run 'mozzy version' to verify the new version"))

		return nil
	},
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func getLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/humancto/mozzy/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var release githubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", err
	}

	// Remove 'v' prefix if present
	version := strings.TrimPrefix(release.TagName, "v")
	return version, nil
}

func detectPlatform() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	platform := goos + "_" + goarch
	return platform
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
