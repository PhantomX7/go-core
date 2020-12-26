package helpers

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// can pass a single url or comma seperated url
// generate s3 object url
func GenerateS3Url(path string) string {
	urls := strings.Split(path, ",")
	for i := range urls {
		if urls[i] == "" {
			if i == 0 {
				return ""
			}
			continue
		}
		urls[i] = url.PathEscape(urls[i])
		urls[i] = fmt.Sprint(
			"https://",
			os.Getenv("AWS_S3_BUCKET"),
			".s3-", os.Getenv("AWS_S3_BUCKET_REGION"),
			".amazonaws.com/",
			urls[i],
		)
	}
	return strings.Join(urls, ",")
}

// can pass a single url or comma seperated url
// generate path of the file from s3Url
func GeneratePathFromS3Url(s3Url string) string {
	urls := strings.Split(s3Url, ",")
	for i := range urls {
		urlFormat := fmt.Sprint(
			"https://",
			os.Getenv("AWS_S3_BUCKET"),
			".s3-", os.Getenv("AWS_S3_BUCKET_REGION"),
			".amazonaws.com/",
		)
		urls[i] = urls[i][len(urlFormat):]
		urls[i], _ = url.PathUnescape(urls[i])
	}
	return strings.Join(urls, ",")
}
