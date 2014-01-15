package version

import (
	"encoding/json"
	"fmt"
	goupdate "github.com/inconshreveable/go-update"
	"how/request"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
)

const (
	Version     = 1.1
	Url         = "https://api.github.com/repos/brettbukowski/how/releases"
	DownloadUrl = "https://github.com/BrettBukowski/how/releases/download/%.1f/how-%s"
)

var releases = make([]map[string]interface{}, 10)
var newestVersion = 0.0

func getReleases() []map[string]interface{} {
	response, _ := request.Get(Url, map[string]string{})

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if err := json.Unmarshal(body, &releases); err != nil {
		panic(err)
	}

	return releases
}

func getNewestVersion() float64 {
	releases := getReleases()
	latestInfo := releases[0]
	latestVersion, _ := strconv.ParseFloat(latestInfo["name"].(string), 64)

	return latestVersion
}

func download(url string) (err error, errRecover error) {
	fmt.Println("Downloading from <%s>", url)

	dl := goupdate.NewDownload(url)

	if err = dl.Get(); err != nil || !dl.Available {
		return
	}

	fmt.Println("Updating...")

	if err, errRecover = goupdate.FromFile(dl.Path); err != nil || errRecover != nil {
		return
	}

	os.Remove(dl.Path)

	fmt.Println("Updated!")

	return
}

func NewerVersionAvailable() bool {
	newest := NewestVersion()

	return newest > Version
}

func NewestVersion() float64 {
	if newestVersion == 0.0 {
		// Cache the operation.
		newestVersion = getNewestVersion()
	}

	return newestVersion
}

func Update() bool {
	if !NewerVersionAvailable() {
		return false
	}

	var err error
	var errRecover error

	if err = goupdate.SanityCheck(); err != nil {
		fmt.Println(err)
		return false
	}

	version := NewestVersion()
	platform := runtime.GOOS

	var target string

	if platform == "darwin" {
		target = "osx"
	} else if platform == "linux" {
		target = "linux-amd64"
	} else if platform == "windows" {
		target = "windows-amd64"
	}

	if err, errRecover = download(fmt.Sprintf(DownloadUrl, version, target)); err != nil || errRecover != nil {
		fmt.Println(err)
		return false
	}

	return true
}
