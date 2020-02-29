package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/bgentry/go-netrc/netrc"
	"github.com/mitchellh/go-homedir"

	terminal "golang.org/x/crypto/ssh/terminal"

	heroku "github.com/heroku/heroku-go/v5"
)

var (
	apiKey      = flag.String("apikey", "", "api key, default found in .netrc, or via username + password")
	createshell = flag.String("createshell", "false", "if missing, the program will spawn a cli")

	appName = flag.String("app", "owlcmsauto", "heroku application to update")
	// archiveName = flag.String("archivename", "owlcms4-heroku", "basename without .tar.gz")

	archive = flag.String("archive", "", "archive url, default is latest inferred from reponame and repoowner")
	// repoName  = flag.String("reponame", *archiveName+"-prerelease", "name of repository")
	// repoOwner = flag.String("repoowner", "jflamy-dev", "owner of repository")

	apiURL = "https://api.heroku.com"
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	spawnCommandWindow()
	err := setupHerokuApiKey(apiKey, apiURL)
	if err != nil {
		defer waitForInput()
		return
	}

	var archiveURL string
	if *archive != "" {
		archiveURL = *archive
		updateApp(appName, "", archiveURL)
	} else {
		updateAllApps()
	}

	waitForInput()
}

func setupHerokuApiKey(apiKey *string, apiURL string) (err error) {
	if *apiKey == "" {
		u, _ := url.Parse(apiURL)
		_, netrcpass, err := getCreds(u)
		if err != nil {
			return err
		}
		heroku.DefaultTransport.Password = netrcpass
	} else {
		heroku.DefaultTransport.Password = *apiKey
	}
	return nil
}

func updateAllApps() {
	h := heroku.NewService(heroku.DefaultClient)
	apps, _ := h.AppList(context.Background(), nil)
	for _, app := range apps {
		appName := app.Name
		configVars, _ := h.ConfigVarInfoForApp(context.Background(), appName)
		latestUrl := configVars["OWLCMS_RELEASES"]
		versionNum := configVars["OWLCMS_VERSION"]
		if latestUrl != nil {
			archiveURL, tagName, err := getArchiveName(*latestUrl)
			if err != nil {
				fmt.Print(err)
			} else {
				ourVersion, _ := semver.NewVersion(tagName)
				theirVersion, _ := semver.NewVersion(*versionNum)
				if ourVersion.GreaterThan(theirVersion) {
					updateApp(&appName, tagName, archiveURL)
				} else {
					fmt.Println(appName + "already up to date (" + *versionNum + ">=" + tagName + ")")
				}

			}
		} else {
			fmt.Println("skipping " + appName)
		}
	}
}

func updateApp(appName *string, tagName string, archiveURL string) {
	h := heroku.NewService(heroku.DefaultClient)
	build, err := h.BuildCreate(context.Background(), *appName, heroku.BuildCreateOpts{
		// anonymous struct. Must include the marshalling tags.
		SourceBlob: struct {
			Checksum *string `json:"checksum,omitempty" url:"checksum,omitempty,key"` // optional tarball checksum
			URL      *string `json:"url,omitempty" url:"url,omitempty,key"`           // URL where gzipped tar archive
			Version  *string `json:"version,omitempty" url:"version,omitempty,key"`   // optional version gzipped tarball.
		}{
			URL: heroku.String(archiveURL),
		}})

	if err != nil {
		log.Print(err)
		defer waitForInput()
		return
	}

	if tagName != "" {
		fmt.Print("updating " + *appName + " to " + tagName + " ")
	} else {
		fmt.Print("updating " + *appName)
	}

	for build.Status == "pending" {
		build, _ = h.BuildInfo(context.Background(), *appName, build.ID)
		fmt.Print(".")
		time.Sleep(time.Second)
	}

	// fix version number
	cviar, err := h.ConfigVarInfoForApp(context.Background(), *appName)
	if err != nil {
		log.Print(err)
		defer waitForInput()
		return
	}
	cviar["OWLCMS_VERSION"] = &tagName
	cviar2, err := h.ConfigVarUpdate(context.Background(), *appName, cviar)
	if err != nil {
		log.Print(err)
		defer waitForInput()
		return
	}
	fmt.Println(" Updated to " + *cviar2["OWLCMS_VERSION"])
}

func waitForInput() {
	if runtime.GOOS == "windows" && *createshell == "spawned" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter any key to close. ")
		_, _ = reader.ReadString('\n')
	}
}

// getCreds obtains the Heroku authorization token
// it tries to get it from the .netrc file
func getCreds(u *url.URL) (user, pass string, err error) {
	netrcPath := getNetRCPath()

	m, err := netrc.FindMachine(netrcPath, u.Host)
	if err != nil {
		// request a token  using username password
		user, token, err := getAuthToken()
		if err != nil {
			// could not authenticate
			return "", "", err
		}
		return user, token, nil
	}

	// .netrc found , check for matching entry
	if m != nil {
		//fmt.Printf("%s - found in netrc (%s): \n", u.Host, netrcPath)
		return m.Login, m.Password, nil
	}
	// request a token using username password
	user, token, err := getAuthToken()
	if err != nil {
		return "", "", err
	}

	//TODO
	return user, token, nil
}

// getUserPassword prompts the user for a Heroku user and password
func getUserPassword() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Heroku Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Heroku Password: ")
	bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	// if err == nil {
	// 	// fmt.Println("\nPassword typed: " + string(bytePassword))
	// }
	password := string(bytePassword)
	fmt.Println()

	return strings.TrimSpace(username), strings.TrimSpace(password)
}

func getAuthToken() (user string, pass string, err error) {
	userName, password := getUserPassword()

	// create request
	client := &http.Client{}
	values := map[string]string{
		"description": "retrieve direct authorization token",
	}
	jsonValue, _ := json.Marshal(values)
	req, _ := http.NewRequest("POST", apiURL+"/oauth/authorizations", bytes.NewBuffer(jsonValue))
	req.SetBasicAuth(userName, password)
	req.Header.Add("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Add("Content-Type", "application/json")

	// parse response
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 201 {
		if resp.StatusCode == 401 {
			fmt.Printf("\nPermission denied. Wrong usename - password combination. (code %v)\n\n", resp.StatusCode)
		} else {
			fmt.Printf("\ncall to API failed (code %v)\n\n", resp.StatusCode)
		}

		if err == nil {
			err = errors.New(strconv.Itoa(resp.StatusCode))
		}
		return "", "", err
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Printf("could not decode %s\n\n", resp.Body)
		return "", "", err
	}

	accessToken := response["access_token"].(map[string]interface{})
	token := accessToken["token"].(string)
	// fmt.Printf("Obtained token %s\n", token)

	_ = addToNetrc(token)

	return "", token, nil
}

func getArchiveName(latestUrl string) (archiveURL string, tagName string, err error) {
	// create request
	client := &http.Client{}
	req, _ := http.NewRequest("GET", latestUrl, nil)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		msg := fmt.Sprintf("call to github %s failed %v", latestUrl, resp.StatusCode)
		fmt.Println(msg)
		defer waitForInput()
		return "", "", errors.New(strconv.Itoa(resp.StatusCode))
	}

	// parse response
	var response []interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		msg := fmt.Sprintf("could not decode %s\n", resp.Body)
		fmt.Println(msg)
		defer waitForInput()
		return "", "", errors.New(strconv.Itoa(resp.StatusCode))
	}

	latest := response[0].(map[string]interface{})
	tagName = latest["tag_name"].(string)
	assets := latest["assets"].([]interface{})
	for _, asset := range assets {
		assetMap := asset.(map[string]interface{})
		archive := assetMap["browser_download_url"].(string)
		//fmt.Printf("Obtained tag name %s - archive = %s \n", tagName, archive)
		return archive, tagName, nil
	}

	msg := fmt.Sprintf("no download found for %s", tagName)
	return "", "", errors.New(msg)
}

func spawnCommandWindow() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	if runtime.GOOS == "windows" && *createshell == "true" {
		// fork a new Command window if running under Windows
		// -createshell false to prevent recursion
		cmd := exec.Command("conhost.exe", ex, "-createshell", "spawned")
		err := cmd.Run()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		os.Exit(0)
	} // else {
	// 	under Linux, assume we are running a command line.
	// }
}

func getNetRCPath() (netrcPath string) {
	hdir, _ := homedir.Dir()
	if runtime.GOOS == "windows" {
		netrcPath = filepath.Join(hdir, "_netrc")
	} else {
		netrcPath = filepath.Join(hdir, ".netrc")
	}
	return netrcPath
}

func addToNetrc(token string) (err error) {
	fd, err := os.OpenFile(getNetRCPath(), os.O_CREATE|os.O_RDWR, os.FileMode(int(0600)))
	if err != nil {
		log.Println(err)
		return err
	}
	defer fd.Close()
	n, err := netrc.Parse(fd)
	if err != nil {
		log.Println(err)
		return err
	}
	n.NewMachine("api.heroku.com", "", token, "")
	text, err := n.MarshalText()
	if err != nil {
		return err
	}
	_, err = fd.WriteAt(text, 0)
	if err != nil {
		return err
	}

	return nil
}
