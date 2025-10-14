package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/humancto/mozzy/internal/formatter"
	"github.com/humancto/mozzy/internal/upload"
)

var (
	uploadFiles      []string
	uploadFieldNames []string
	uploadFields     []string
	uploadNoProgress bool
)

var uploadCmd = &cobra.Command{
	Use:   "upload <url>",
	Short: "Upload files with multipart form data",
	Long: `Upload files to a server using multipart/form-data.

Supports:
  - Single or multiple file uploads
  - Custom form fields
  - Custom headers and authentication
  - Progress tracking

Examples:
  # Upload single file
  mozzy upload https://api.example.com/upload -f avatar.jpg

  # Upload with custom field name
  mozzy upload https://api.example.com/upload -f avatar.jpg --field-name profileImage

  # Upload multiple files
  mozzy upload https://api.example.com/upload -f file1.jpg -f file2.png

  # Upload with form fields
  mozzy upload https://api.example.com/upload -f resume.pdf --data "name=John" --data "email=john@example.com"

  # Upload with authentication
  mozzy upload https://api.example.com/upload -f file.jpg --auth token123

  # Upload without progress
  mozzy upload https://api.example.com/upload -f file.jpg --no-progress`,
	Args: cobra.ExactArgs(1),
	RunE: runUpload,
}

func init() {
	uploadCmd.Flags().StringSliceVarP(&uploadFiles, "file", "f", nil, "File to upload (repeatable)")
	uploadCmd.Flags().StringSliceVar(&uploadFieldNames, "field-name", nil, "Form field name for file (repeatable, matches --file order)")
	uploadCmd.Flags().StringSliceVar(&uploadFields, "data", nil, "Form field (format: name=value, repeatable)")
	uploadCmd.Flags().BoolVar(&uploadNoProgress, "no-progress", false, "Disable progress display")
	uploadCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(uploadCmd)
}

func runUpload(cmd *cobra.Command, args []string) error {
	url := args[0]

	if len(uploadFiles) == 0 {
		return fmt.Errorf("at least one file is required (use -f or --file)")
	}

	// Build file uploads
	files := make([]upload.FileUpload, 0, len(uploadFiles))
	for i, filePath := range uploadFiles {
		fieldName := "file"
		if i < len(uploadFieldNames) && uploadFieldNames[i] != "" {
			fieldName = uploadFieldNames[i]
		} else if len(uploadFiles) > 1 {
			fieldName = fmt.Sprintf("file%d", i)
		}

		files = append(files, upload.FileUpload{
			FieldName: fieldName,
			FilePath:  filePath,
		})
	}

	// Parse form fields
	fields := make([]upload.FormField, 0, len(uploadFields))
	for _, field := range uploadFields {
		parts := strings.SplitN(field, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid field format: %q (expected name=value)", field)
		}
		fields = append(fields, upload.FormField{
			Name:  parts[0],
			Value: parts[1],
		})
	}

	// Build headers map
	headersMap := make(map[string]string)
	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			headersMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	// Upload
	opts := upload.UploadOptions{
		URL:          url,
		Files:        files,
		Fields:       fields,
		Headers:      headersMap,
		AuthToken:    authToken,
		ShowProgress: !uploadNoProgress,
		Timeout:      30 * time.Minute,
	}

	resp, body, err := upload.UploadWithProgress(opts)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	defer resp.Body.Close()

	// Print response
	formatter.PrintStatusLine("POST", url, resp.StatusCode, 0)
	if err := formatter.PrintJSONOrText(body, jqQuery); err != nil {
		return err
	}

	if failOnErr && resp.StatusCode >= 400 {
		return fmt.Errorf("upload failed with status: %s", resp.Status)
	}

	return nil
}
