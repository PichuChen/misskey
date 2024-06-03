package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/misskey-dev/misskey/models/user"
)

var misskeyNodeBackend = "http://localhost:3000"
var backend2ListenPort = "3001"

func main() {
	user.InitDB()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Redirdect to port 3000
		if r.URL.Path == "/api/users/show" {
			routeApiUsersShow(w, r)
			return
		}
		proxyToNodeBackend(r, w)

	})
	http.ListenAndServe(":"+backend2ListenPort, mux)

}

func routeApiUsersShow(w http.ResponseWriter, r *http.Request) {

	// decode json body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Error reading request body", "err", err)
		return
	}
	slog.Info("Request body", "body", string(b))
	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		slog.Error("Error decoding json", "err", err)
		return
	}
	slog.Info("Decoded json", "json", m)
	// proxyToNodeBackendWithPayload(r, w, b)
	// return
	if m["username"] == nil {
		slog.Error("Username not provided")
		return
	}

	u, err := user.FindOneByUsername(m["username"].(string))
	if err != nil {
		slog.Error("Error finding user", "err", err)
		return
	}
	if u == nil {
		slog.Info("User not found")
		return
	}
	// slog.Info("User found", "user", u)

	o := packMiUser(u)
	// slog.Info("Packed user", "user", o)
	j, err := json.Marshal(o)
	if err != nil {
		slog.Error("Error marshalling user", "err", err)
		return
	}

	w.Write(j)

}

func getIdenticonUrl(user *user.MiUser) string {
	// return `${this.config.url}/identicon/${user.username.toLowerCase()}@${user.host ?? this.config.host}`;
	if user.Host == nil {
		return fmt.Sprintf("%v/identicon/%s@%v", misskeyNodeBackend, strings.ToLower(user.Username), "127.0.0.1:3000")
	}
	return fmt.Sprintf("%v/identicon/%s@%v", misskeyNodeBackend, strings.ToLower(user.Username), user.Host)
}

func packMiUser(u *user.MiUser) map[string]interface{} {
	o := make(map[string]interface{})
	o["id"] = u.ID
	o["updatedAt"] = u.UpdatedAt
	o["lastFetchedAt"] = u.LastFetchedAt
	o["lastActiveDate"] = u.LastActiveDate
	o["hideOnlineStatus"] = u.HideOnlineStatus
	o["username"] = u.Username
	o["usernameLower"] = u.UsernameLower
	o["name"] = u.Name
	o["followersCount"] = u.FollowersCount
	o["followingCount"] = u.FollowingCount
	o["movedToUri"] = u.MovedToURI
	o["movedAt"] = u.MovedAt
	o["alsoKnownAs"] = u.AlsoKnownAs
	o["notesCount"] = u.NotesCount
	o["avatarId"] = u.AvatarID
	o["bannerId"] = u.BannerID
	if u.AvatarURL == nil {
		o["avatarUrl"] = getIdenticonUrl(u)
	} else {
		o["avatarUrl"] = u.AvatarURL
	}
	o["bannerUrl"] = u.BannerURL
	o["avatarBlurhash"] = u.AvatarBlurHash
	o["bannerBlurhash"] = u.BannerBlurHash
	o["avatarDecorations"] = u.AvatarDecorations
	o["tags"] = u.Tags
	o["isSuspended"] = u.IsSuspended
	o["isLocked"] = u.IsLocked
	o["isBot"] = u.IsBot
	o["isCat"] = u.IsCat
	o["isRoot"] = u.IsRoot
	o["isExplorable"] = u.IsExplorable
	o["isHibernated"] = u.IsHibernated
	o["isDeleted"] = u.IsDeleted
	o["emojis"] = u.Emojis
	o["host"] = u.Host
	o["inbox"] = u.Inbox
	o["sharedInbox"] = u.SharedInbox
	o["featured"] = u.Featured
	o["uri"] = u.URI
	o["followersUri"] = u.FollowersURI
	o["token"] = u.Token
	return o

}

func proxyToNodeBackend(r *http.Request, w http.ResponseWriter) {

	requestURL := r.URL
	requestMethod := r.Method
	requestHeaders := r.Header
	slog.Info("Request received", "method", requestMethod, "url", requestURL, "headers", requestHeaders)

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Error reading request body", "err", err)
		return
	}
	slog.Info("Request body", "body", string(requestBody))

	proxyToNodeBackendWithPayload(r, w, requestBody)
}
func proxyToNodeBackendWithPayload(r *http.Request, w http.ResponseWriter, requestBody []byte) {
	requestURL := r.URL
	requestMethod := r.Method
	requestHeaders := r.Header
	client := &http.Client{}
	req, err := http.NewRequest(requestMethod, misskeyNodeBackend+requestURL.String(), bytes.NewReader(requestBody))
	if err != nil {
		slog.Error("Error creating request", "err", err)
		return
	}

	req.Header = requestHeaders
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error sending request", "err", err)
		return
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body", "err", err)
		return
	}

	for key, value := range resp.Header {
		w.Header().Set(key, value[0])
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)

	slog.Info("Request sent")

}
