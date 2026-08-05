package lfm

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"pymouse/pymouse/config"
	"strings"
	"time"

	"github.com/google/uuid"
)

// r2Configured reports whether all required R2 env vars are set. When false,
// callers should fall back to a non-image inline result.
func r2Configured() bool {
	return config.R2AccountID != "" && config.R2Bucket != "" &&
		config.R2AccessKeyID != "" && config.R2SecretKey != "" && config.R2PublicURL != ""
}

// r2MissingEnv lists which R2_* env vars are unset, for clearer error logs.
func r2MissingEnv() []string {
	var missing []string
	if config.R2AccountID == "" {
		missing = append(missing, "R2_ACCOUNT_ID")
	}
	if config.R2Bucket == "" {
		missing = append(missing, "R2_BUCKET")
	}
	if config.R2AccessKeyID == "" {
		missing = append(missing, "R2_ACCESS_KEY_ID")
	}
	if config.R2SecretKey == "" {
		missing = append(missing, "R2_SECRET_ACCESS_KEY")
	}
	if config.R2PublicURL == "" {
		missing = append(missing, "R2_PUBLIC_URL")
	}
	return missing
}

// uploadToR2 uploads the JPEG at filePath to the configured Cloudflare R2
// bucket using a presigned PUT (AWS Signature Version 4), and returns the
// object's public URL. No AWS SDK dependency — the signature is built with
// the standard library. The object key is a random UUID so concurrent /
// repeated uploads never collide.
func uploadToR2(filePath string) (string, error) {
	if !r2Configured() {
		return "", fmt.Errorf("R2 not configured (missing: %v)", r2MissingEnv())
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	key := "lfm-inline/" + uuid.NewString() + ".jpg"
	host := fmt.Sprintf("%s.r2.cloudflarestorage.com", config.R2AccountID)
	endpoint := fmt.Sprintf("https://%s/%s/%s", host, config.R2Bucket, key)
	publicURL := strings.TrimRight(config.R2PublicURL, "/") + "/" + key

	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "image/jpeg")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	req.Header.Set("Host", host)

	if err := signS3V4PUT(req, key, data); err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return "", fmt.Errorf("R2 upload: HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return publicURL, nil
}

// signS3V4PUT applies an AWS Signature Version 4 Authorization header to req
// for a PUT to the configured R2 bucket. Suffix "r2" maps to region "auto".
func signS3V4PUT(req *http.Request, objectKey string, payload []byte) error {
	const (
		region    = "auto"
		service   = "s3"
		algorithm = "AWS4-HMAC-SHA256"
	)

	now := time.Now().UTC()
	dateStamp := now.Format("20060102")
	amzDate := now.Format("20060102T150405Z")
	payloadHash := sha256Hex(payload)

	// S3 requires these for SigV4.
	req.Header.Set("x-amz-date", amzDate)
	req.Header.Set("x-amz-content-sha256", payloadHash)
	req.Header.Set("x-amz-storage-class", "STANDARD")

	// 1. Canonical request.
	signedHeaders := "host;x-amz-content-sha256;x-amz-date;x-amz-storage-class"
	canonicalHeaders := strings.Join([]string{
		req.Host + "\n",
		payloadHash + "\n",
		amzDate + "\n",
		"STANDARD\n",
	}, "")
	canonicalURI := "/" + config.R2Bucket + "/" + objectKey
	canonicalRequest := strings.Join([]string{
		"PUT",
		canonicalURI,
		"", // no query string
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	}, "\n")

	// 2. String to sign.
	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", dateStamp, region, service)
	stringToSign := strings.Join([]string{
		algorithm,
		amzDate,
		credentialScope,
		sha256Hex([]byte(canonicalRequest)),
	}, "\n")

	// 3. Signing key.
	kDate := hmacSHA256([]byte("AWS4"+config.R2SecretKey), []byte(dateStamp))
	kRegion := hmacSHA256(kDate, []byte(region))
	kService := hmacSHA256(kRegion, []byte(service))
	kSigning := hmacSHA256(kService, []byte("aws4_request"))

	// 4. Signature.
	signature := hex.EncodeToString(hmacSHA256(kSigning, []byte(stringToSign)))

	authz := fmt.Sprintf(
		"%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		algorithm, config.R2AccessKeyID, credentialScope, signedHeaders, signature,
	)
	req.Header.Set("Authorization", authz)
	return nil
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func hmacSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
