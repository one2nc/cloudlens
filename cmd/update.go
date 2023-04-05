package cmd

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/charmbracelet/glamour"
	"github.com/cheggaaa/pb/v3"
	"github.com/google/go-github/v50/github"
	"github.com/minio/selfupdate"
	"github.com/spf13/cobra"
)

type (
	AssetId     int
	AssetFormat uint
)

const (
	Zip AssetFormat = iota
	Tar
)

var (
	Oraganisation         = "one2nc"
	ToolName              = "cloudlens"
	extIfFound            = ".exe"
	HideProgressBar       = false
	HideReleaseNotes      = false
	DownloadUpdateTimeout = time.Duration(30) * time.Second

	httpClient = &http.Client{
		Timeout: DownloadUpdateTimeout,
	}
	githubClient = github.NewClient(httpClient)
)

func updateCmd() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Updates cloudlens to the latest version",
		Run: func(cmd *cobra.Command, args []string) {
			if err := doUpdate(); err != nil {
				fmt.Printf("%v\n", err.Error())
			}
		},
	}
	return updateCmd
}

func doUpdate() error {
	latest, _, err := githubClient.Repositories.GetLatestRelease(context.Background(), Oraganisation, ToolName)
	if err != nil {
		return err
	}
	latestVersion, err := semver.NewVersion(latest.GetTagName())
	if err != nil {
		return fmt.Errorf("[ERROR] failed to parse semversion from tagname `%v` got %v", latest.GetTagName(), err)
	}
	currentVersion, err := semver.NewVersion(version)
	if err != nil {
		return fmt.Errorf("[ERROR] failed to parse semversion from current version %v got %v", version, err)
	}
	if !latestVersion.GreaterThan(currentVersion) {
		return fmt.Errorf("%v is already updated to latest version", ToolName)
	}
	assetId, format, err := getAssetIDAndFormatFromRelease(latest)
	if err != nil {
		return err
	}
	bin, err := getExecutableFromAsset(assetId, format)
	if err != nil {
		return fmt.Errorf("[ERROR] executable %v not found in release asset `%v` got: %v", ToolName, assetId, err)
	}

	updateOpts := selfupdate.Options{}
	if err = selfupdate.Apply(bytes.NewBuffer(bin), updateOpts); err != nil {
		fmt.Printf("Error] update of %v %v -> %v failed, rolling back update\n", ToolName, currentVersion.String(), latestVersion.String())
		if err := selfupdate.RollbackError(err); err != nil {
			fmt.Printf("rollback of update of %v failed got %v,pls reinstall %v\n", ToolName, err, ToolName)
		}
		os.Exit(0)
	}

	fmt.Printf("%v sucessfully updated %v -> %v (latest)\n", ToolName, currentVersion.String(), latestVersion.String())

	if !HideReleaseNotes {
		output := latest.GetBody()
		if rendered, err := glamour.Render(output, "dark"); err == nil {
			output = rendered
		} else {
			fmt.Printf("[Error] %v\n", err.Error())
		}
		fmt.Printf("%v\n\n", output)
	}
	os.Exit(0)
	return err
}

// getAssetIDAndFormatFromRelease finds AssetID and formatfrom release or returns a descriptive error
func getAssetIDAndFormatFromRelease(latest *github.RepositoryRelease) (AssetId, AssetFormat, error) {
	builder := &strings.Builder{}
	builder.WriteString("cloudlens")
	builder.WriteString("_")
	builder.WriteString(strings.TrimPrefix(latest.GetTagName(), "v"))
	builder.WriteString("_")
	if strings.EqualFold(runtime.GOOS, "darwin") {
		builder.WriteString("darwin_all")
	} else {
		builder.WriteString(runtime.GOOS)
		builder.WriteString("_")
		builder.WriteString(runtime.GOARCH)
	}

	var assetId AssetId
	var format AssetFormat
loop:
	for _, v := range latest.Assets {
		asset := *v.Name
		switch {
		case strings.Contains(asset, ".zip"):
			if strings.EqualFold(asset, builder.String()+".zip") {
				assetId = AssetId(*v.ID)
				format = Zip
				break loop
			}
		case strings.Contains(asset, ".tar.gz"):
			if strings.EqualFold(asset, builder.String()+".tar.gz") {
				assetId = AssetId(*v.ID)
				format = Tar
				break loop
			}
		}
	}
	builder.Reset()

	if assetId == 0 {
		return assetId, format, fmt.Errorf("%v", "Asset not found!")
	}
	return assetId, format, nil
}

func getExecutableFromAsset(assetId AssetId, format AssetFormat) ([]byte, error) {
	buff, err := downloadAssetFromID(assetId)
	if err != nil {
		return nil, err
	}
	if format == Zip {
		zipReader, err := zip.NewReader(bytes.NewReader(buff.Bytes()), int64(buff.Len()))
		if err != nil {
			return nil, err
		}
		for _, f := range zipReader.File {
			if !strings.EqualFold(strings.TrimSuffix(f.Name, extIfFound), ToolName) {
				continue
			}
			fileInArchive, err := f.Open()
			if err != nil {
				return nil, err
			}
			bin, err := io.ReadAll(fileInArchive)
			if err != nil {
				return nil, err
			}
			_ = fileInArchive.Close()
			return bin, nil
		}
	} else if format == Tar {
		gzipReader, err := gzip.NewReader(buff)
		if err != nil {
			return nil, err
		}
		tarReader := tar.NewReader(gzipReader)
		// iterate through the files in the archive
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			if !strings.EqualFold(strings.TrimSuffix(header.FileInfo().Name(), extIfFound), ToolName) {
				continue
			}
			// if the file is not a directory, extract it
			if !header.FileInfo().IsDir() {
				bin, err := io.ReadAll(tarReader)
				if err != nil {
					return nil, err
				}
				return bin, nil
			}
		}
	}
	return nil, fmt.Errorf("executable not found in archive")
}

// downloadAssetFromID downloads and returns a buffer or a descriptive error
func downloadAssetFromID(assetId AssetId) (*bytes.Buffer, error) {
	_, rdurl, err := githubClient.Repositories.DownloadReleaseAsset(context.Background(), Oraganisation, ToolName, int64(assetId), nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Get(rdurl)
	if err != nil {
		return nil, fmt.Errorf("%v: failed to downlaod asset", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("something went wrong got %v while downloading asset, expected status 200", resp.StatusCode)
	}
	if resp.Body == nil {
		return nil, errors.New("something went wrong got response without body")
	}
	defer resp.Body.Close()

	if !HideProgressBar {
		bar := pb.New64(resp.ContentLength).SetMaxWidth(100)
		bar.Start()
		resp.Body = bar.NewProxyReader(resp.Body)
		defer bar.Finish()
	}

	bin, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%v :failed to read response body", err)
	}
	return bytes.NewBuffer(bin), nil
}
