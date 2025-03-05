package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var _ function.Function = &DownloadFileFunction{}

type DownloadFileFunction struct {
}

func NewDownloadFileFunction() function.Function {
	return &DownloadFileFunction{}
}

func (d *DownloadFileFunction) Metadata(ctx context.Context, request function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "file"
}

func (d *DownloadFileFunction) Definition(ctx context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Summary:     "Downloads a file, returning the filename.",
		Description: "Downloads a file from a given URL and returns the filename.",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "url",
				Description: "URL to download",
			},
			function.StringParameter{
				Name:        "filename",
				Description: "Name of the filename for the contents.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (d *DownloadFileFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var url string
	var filename string
	skipDownload := false

	response.Error = function.ConcatFuncErrors(response.Error, request.Arguments.Get(ctx, &url, &filename))

	if !isValidURL(url) {
		response.Error = function.NewFuncError("invalid url")
		return
	}

	if filename == "" {
		response.Error = function.NewFuncError("filename is empty")
		return
	}

	info, _ := os.Stat(filename)
	if info != nil {
		_, contentLength, err := getRemoteFileMetadata(url)
		if err != nil {
			response.Error = function.NewFuncError(fmt.Sprintf("error getting remote metadata: %v", err))
		}

		if info.Size() == contentLength {
			skipDownload = true
		}
	}

	if !skipDownload {
		err := downloadFile(filename, url)
		if err != nil {
			response.Error = function.NewFuncError(fmt.Sprintf("error downloading file: %v", err))
			return
		}
	}

	response.Error = function.ConcatFuncErrors(response.Error, response.Result.Set(ctx, filename))
}

func getRemoteFileMetadata(url string) (etag string, contentLength int64, err error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("error closing response body: %s", err)
		}
	}(resp.Body)

	etag = resp.Header.Get("ETag")

	contentLengthStr := resp.Header.Get("Content-Length")
	if contentLengthStr != "" {
		contentLength, err = strconv.ParseInt(contentLengthStr, 10, 64)
		if err != nil {
			return "", 0, err
		}
	}

	return etag, contentLength, nil
}
