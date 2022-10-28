package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/luigizuccarelli/golang-gitwebhook-service/pkg/connectors"
	"github.com/luigizuccarelli/golang-gitwebhook-service/pkg/schema"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
	ERRMSG          string = "{\"status\":\"KO\", \"statuscode\":\"500\",\"message\":\""
)

func WebhookHandler(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	var git *schema.GitSchema
	var mapping *schema.MapBinding
	var eventListenerUrl string
	var payload string

	body, err := ioutil.ReadAll(r.Body)
	if strings.Contains(string(body), "payload=") {
		formatted := strings.Split(string(body), "=")[1]
		payload, err = url.QueryUnescape(formatted)
		if err != nil {
			log.Fatal(err)
			return
		}
	} else {
		payload = string(body)
	}

	con.Trace("Input data %s", payload)
	if err != nil {
		con.Error("WebhookHandler could not read body data %v", err)
		resp := ERRMSG + fmt.Sprintf("\"WebhookHandler could not read body data %v", err) + "\"}"
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", resp)
		return
	}

	err = json.Unmarshal([]byte(payload), &git)
	if err != nil {
		con.Error("WebhookHandler could not unmarshal to struct %v", err)
		resp := ERRMSG + fmt.Sprintf("\"WebhookHandler could not unmarshal struct %v", err) + "\"}"
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", resp)
		return
	}

	con.Debug("Mapping struct %v", git)

	postRequest := false
	// post to the various eventlisteners
	if len(os.Getenv("PR_OPENED_URL")) > 0 && git.Action == "opened" {
		// post to the pr_opened eventlistener
		mapping = &schema.MapBinding{
			RepoUrl:   git.Repository.CloneURL,
			RepoName:  git.Repository.Name,
			RepoHash:  git.PullRequest.Head.Sha,
			ActorName: git.PullRequest.User.Login,
			Message:   git.PullRequest.Title,
		}
		postRequest = true
		eventListenerUrl = os.Getenv("PR_OPENED_URL")
	}

	if len(os.Getenv("PR_MERGED_URL")) > 0 && git.Action == "closed" {
		// only post on merged true
		if git.PullRequest.Merged {
			// post to the pr+merged eventlistener
			mapping = &schema.MapBinding{
				RepoUrl:   git.Repository.CloneURL,
				RepoName:  git.Repository.Name,
				RepoHash:  git.PullRequest.MergeCommitSha,
				ActorName: git.PullRequest.User.Login,
				Message:   git.PullRequest.Title,
			}
			postRequest = true
			eventListenerUrl = os.Getenv("PR_MERGED_URL")
		}
	}

	if (len(os.Getenv("PRERELEASED_URL")) > 0 || len(os.Getenv("RELEASED_URL")) > 0) && git.Action == "published" {
		// check for the prerelease field
		if git.Release.Prerelease {
			// post to the pre_release eventlistener
			mapping = &schema.MapBinding{
				RepoUrl:    git.Repository.CloneURL,
				RepoName:   git.Repository.Name,
				RepoHash:   git.Release.TargetCommitish,
				ActorName:  git.Release.Author.Login,
				Message:    git.Release.Name + " " + git.Release.Body,
				TagVersion: git.Release.TagName,
			}
			postRequest = true
			eventListenerUrl = os.Getenv("PRERELEASED_URL")
		} else {
			// post to release eventlistener
			mapping = &schema.MapBinding{
				RepoUrl:    git.Repository.CloneURL,
				RepoName:   git.Repository.Name,
				RepoHash:   git.Release.TargetCommitish,
				ActorName:  git.Release.Author.Login,
				Message:    git.Release.Name + " " + git.Release.Body,
				TagVersion: git.Release.TagName,
			}
			postRequest = true
			eventListenerUrl = os.Getenv("RELEASED_URL")
		}
	}

	if postRequest {
		_, err = makePostRequest(eventListenerUrl, APPLICATIONJSON, mapping, con)
		if err != nil {
			resp := ERRMSG + fmt.Sprintf("\"Request failed %v", err) + "\"}"
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", resp)
			return
		}
		resp := "{\"status\":\"OK\", \"statuscode\":\"200\",\"message\":\"Request sent successfully\"}"
		w.WriteHeader(http.StatusOK)
		con.Debug("Result struct for git webhook %v", mapping)
		fmt.Fprintf(w, "%s", string(resp))
	} else {
		con.Info("NOP (opened,merged,release or prerelease action not detected)")
	}
}

func IsAlive(w http.ResponseWriter, r *http.Request, con connectors.Clients) {
	con.Trace("Request Object", r)
	fmt.Fprintf(w, "%s", "{\"name\":\"golang-gitwebhook-service\",\"version\":\"v0.0.1\"}")
}

// makePostRequest - private utility function for POST
func makePostRequest(elUrl string, contentType string, mb *schema.MapBinding, con connectors.Clients) ([]byte, error) {
	var b []byte

	data, _ := json.MarshalIndent(mb, "", "    ")
	req, _ := http.NewRequest("POST", elUrl, bytes.NewBuffer(data))
	con.Debug("Post data to eventListenerUrl : %s", string(data))
	req.Header.Set(CONTENTTYPE, contentType)
	con.Info("Function makeRequest %s", elUrl)
	resp, err := con.Do(req)
	if err != nil {
		con.Error("Function makePostRequest http request %v", err)
		return b, err
	}
	defer resp.Body.Close()
	if resp.StatusCode <= http.StatusAccepted {
		con.Debug("Function makePostRequest response from middleware %d", resp.StatusCode)
		return []byte("ok"), nil
	}
	con.Error("Function makePostRequest response code %v", resp.StatusCode)
	return []byte("ko"), errors.New(strconv.Itoa(resp.StatusCode))
}
