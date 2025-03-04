package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"io"
	"log"
	"net/http"
	"os"
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

	response.Error = function.ConcatFuncErrors(response.Error, request.Arguments.Get(ctx, &url, &filename))

	err := downloadFileFunc(filename, url)
	if err != nil {
		response.Error = function.NewFuncError(fmt.Sprintf("error downloading file: %v", err))
	}

	response.Error = function.ConcatFuncErrors(response.Error, response.Result.Set(ctx, filename))
}

func downloadFileFunc(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("error closing response body: %s", err)
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		err := out.Close()
		if err != nil {
			log.Printf("error closing file output: %s", err)
		}
	}()

	_, err = io.Copy(out, resp.Body)
	return nil
}
