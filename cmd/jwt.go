package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/humancto/mozzy/internal/jwtutil"
)

var (
	jwtSecret string
	jwkURL    string
	alg       string
)

func init() {
	jwtCmd.PersistentFlags().StringVar(&jwtSecret, "secret", "", "HMAC secret (HS256/384/512)")
	jwtCmd.PersistentFlags().StringVar(&jwkURL, "jwk", "", "JWKS URL for RSA/ECDSA verification")
	jwtCmd.PersistentFlags().StringVar(&alg, "alg", "HS256", "Signing algorithm for `sign` (HS256/384/512)")
	rootCmd.AddCommand(jwtCmd)
}

var jwtCmd = &cobra.Command{ Use: "jwt", Short: "Decode, verify, and sign JWTs" }

var jwtDecode = &cobra.Command{
	Use:   "decode <token>",
	Short: "Decode a JWT without verifying",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		h, p, err := jwtutil.Decode(args[0])
		if err != nil { return err }
		fmt.Println("Header:")
		for k, v := range h { fmt.Printf("  %s: %v\n", k, v) }
		fmt.Println("\nPayload:")
		for k, v := range p { fmt.Printf("  %s: %v\n", k, v) }
		if exp, ok := p["exp"].(float64); ok {
			t := time.Unix(int64(exp), 0)
			fmt.Printf("\nexp (human): %s (%s)\n", t.Format(time.RFC3339), time.Until(t).Round(time.Second))
		}
		return nil
	},
}

var jwtVerify = &cobra.Command{
	Use:   "verify <token>",
	Short: "Verify a JWT with HMAC secret or JWKS URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		switch {
		case jwkURL != "":
			err = jwtutil.VerifyWithJWKS(args[0], jwkURL)
		case jwtSecret != "":
			err = jwtutil.VerifyHMAC(args[0], []byte(jwtSecret))
		default:
			return fmt.Errorf("provide --secret (HS*) or --jwk (JWKS)")
		}

		if err != nil {
			return fmt.Errorf("❌ verification failed: %w", err)
		}

		fmt.Println("✅ JWT signature is valid")

		// Also decode and show expiration if present
		_, p, _ := jwtutil.Decode(args[0])
		if exp, ok := p["exp"].(float64); ok {
			t := time.Unix(int64(exp), 0)
			remaining := time.Until(t)
			if remaining > 0 {
				fmt.Printf("⏰ Expires in: %s (at %s)\n", remaining.Round(time.Second), t.Format(time.RFC3339))
			} else {
				fmt.Printf("⚠️  Token expired %s ago (at %s)\n", (-remaining).Round(time.Second), t.Format(time.RFC3339))
			}
		}

		return nil
	},
}

var jwtSign = &cobra.Command{
	Use:   "sign <payload.json>",
	Short: "Sign a JSON payload (HS* family)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if jwtSecret == "" {
			return fmt.Errorf("--secret required for signing (HS*)")
		}
		b, err := os.ReadFile(args[0])
		if err != nil { return err }
		tok, err := jwtutil.SignHMAC(b, []byte(jwtSecret), alg)
		if err != nil { return err }
		fmt.Println(tok)
		return nil
	},
}

func init() { jwtCmd.AddCommand(jwtDecode, jwtVerify, jwtSign) }
