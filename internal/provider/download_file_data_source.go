package provider

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io"
	"net/http"
	"os"
)

var _ datasource.DataSource = &DownloadFileDataSource{}

type DownloadFileDataSource struct {
}

func NewDownloadFileDataSource() datasource.DataSource {
	return &DownloadFileDataSource{}
}

type DownloadFileDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	Url          types.String `tfsdk:"url"`
	OutputFile   types.String `tfsdk:"output_file"`
	Base64SHA256 types.String `tfsdk:"output_base64sha256"`
	MD5          types.String `tfsdk:"output_md5"`
	SHA          types.String `tfsdk:"output_sha"`
	SHA256       types.String `tfsdk:"output_sha256"`
	FileSize     types.Int64  `tfsdk:"output_size"`
	VerifySHA256 types.String `tfsdk:"verify_sha256"`
	VerifySHA    types.String `tfsdk:"verify_sha"`
	VerifyMD5    types.String `tfsdk:"verify_md5"`
}

func (f *DownloadFileDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_file"
}

func (f *DownloadFileDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "Downloads a file from a website using the supplied URL.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "URL to download",
				Required:            true,
			},
			"output_file": schema.StringAttribute{
				MarkdownDescription: "File name to write content",
				Required:            true,
			},
			"output_base64sha256": schema.StringAttribute{
				MarkdownDescription: "Base64 Encoded SHA256 checksum of output file",
				Computed:            true,
			},
			"output_md5": schema.StringAttribute{
				MarkdownDescription: "MD5 of output file",
				Computed:            true,
			},
			"output_sha": schema.StringAttribute{
				MarkdownDescription: "SHA1 checksum of output file",
				Computed:            true,
			},
			"output_sha256": schema.StringAttribute{
				MarkdownDescription: "SHA256 checksum of output file",
				Computed:            true,
			},
			"output_size": schema.Int64Attribute{
				MarkdownDescription: "File size of output file",
				Computed:            true,
			},
			"verify_sha256": schema.StringAttribute{
				MarkdownDescription: "SHA256 checksum to verify",
				Optional:            true,
			},
			"verify_sha": schema.StringAttribute{
				MarkdownDescription: "SHA1 checksum to verify",
				Optional:            true,
			},
			"verify_md5": schema.StringAttribute{
				MarkdownDescription: "MD5 checksum to verify",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
		},
	}
}

func (f *DownloadFileDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data DownloadFileDataSourceModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	err := downloadFile(data.OutputFile.ValueString(), data.Url.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Download file error", err.Error())
		return
	}

	fi, err := os.Stat(data.OutputFile.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Download file error", err.Error())
		return
	}

	data.FileSize = types.Int64Value(fi.Size())

	err = genFileShas(data.OutputFile.ValueString(), &data)
	if err != nil {
		response.Diagnostics.AddError("Download file error", err.Error())
		return
	}

	err = verifyFileShas(&data)
	if err != nil {
		response.Diagnostics.AddError("Download file error", err.Error())
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func genFileShas(filename string, data *DownloadFileDataSourceModel) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not compute file '%s' checksum: %s", filename, err)
	}

	h := sha1.New()
	h.Write(content)
	sha1Hash := hex.EncodeToString(h.Sum(nil))

	h256 := sha256.New()
	h256.Write(content)
	shaSum := h256.Sum(nil)
	sha256Hash := hex.EncodeToString(h256.Sum(nil))
	sha256base64 := base64.StdEncoding.EncodeToString(shaSum[:])

	md5Hash := md5.New()
	md5Hash.Write(content)
	md5Sum := hex.EncodeToString(md5Hash.Sum(nil))

	data.SHA = types.StringValue(sha1Hash)
	data.SHA256 = types.StringValue(sha256Hash)
	data.Base64SHA256 = types.StringValue(sha256base64)
	data.MD5 = types.StringValue(md5Sum)
	data.Id = types.StringValue(sha1Hash)

	return nil
}

func verifyFileShas(data *DownloadFileDataSourceModel) error {
	if !data.VerifySHA256.IsNull() {
		if data.VerifySHA256.ValueString() != data.SHA256.ValueString() {
			return errors.New("SHA256 signature mismatch")
		}
	}

	if !data.VerifySHA.IsNull() {
		if data.VerifySHA.ValueString() != data.SHA.ValueString() {
			return errors.New("SHA1 signature mismatch")
		}
	}

	if !data.VerifyMD5.IsNull() {
		if data.VerifyMD5.ValueString() != data.MD5.ValueString() {
			return errors.New("MD5 signature mismatch")
		}
	}

	return nil
}
