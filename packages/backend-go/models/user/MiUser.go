package user

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

type MiUser struct {
	ID               string     `json:"id"`
	UpdatedAt        *time.Time `json:"updatedAt"`
	LastFetchedAt    *time.Time `json:"lastFetchedAt"`
	LastActiveDate   *time.Time `json:"lastActiveDate"`
	HideOnlineStatus bool       `json:"hideOnlineStatus"`
	Username         string     `json:"username"`
	UsernameLower    string     `json:"usernameLower"`
	Name             *string    `json:"name"`
	FollowersCount   int        `json:"followersCount"`
	FollowingCount   int        `json:"followingCount"`
	MovedToURI       *string    `json:"movedToUri"`
	MovedAt          *time.Time `json:"movedAt"`
	AlsoKnownAs      *[]string  `json:"alsoKnownAs"`
	NotesCount       int        `json:"notesCount"`
	AvatarID         *string    `json:"avatarId"`
	// Avatar 			 *MiDriveFile     `json:"avatar"`
	BannerID *string `json:"bannerId"`
	// Banner 			 *MiDriveFile     `json:"banner"`
	AvatarURL         *string `json:"avatarUrl"`
	BannerURL         *string `json:"bannerUrl"`
	AvatarBlurHash    *string `json:"avatarBlurhash"`
	BannerBlurHash    *string `json:"bannerBlurhash"`
	AvatarDecorations []struct {
		ID      string `json:"id"`
		Angle   int    `json:"angle"`
		FlipH   bool   `json:"flipH"`
		OffsetX int    `json:"offsetX"`
		OffsetY int    `json:"offsetY"`
	} `json:"avatarDecorations"`
	Tags         []string `json:"tags"`
	IsSuspended  bool     `json:"isSuspended"`
	IsLocked     bool     `json:"isLocked"`
	IsBot        bool     `json:"isBot"`
	IsCat        bool     `json:"isCat"`
	IsRoot       bool     `json:"isRoot"`
	IsExplorable bool     `json:"isExplorable"`
	IsHibernated bool     `json:"isHibernated"`
	// アカウントが削除されたかどうかのフラグだが、完全に削除される際は物理削除なので実質削除されるまでの「削除が進行しているかどうか」のフラグ
	IsDeleted    bool     `json:"isDeleted"`
	Emojis       []string `json:"emojis"`
	Host         *string  `json:"host"`
	Inbox        *string  `json:"inbox"`
	SharedInbox  *string  `json:"sharedInbox"`
	Featured     *string  `json:"featured"`
	URI          *string  `json:"uri"`
	FollowersURI *string  `json:"followersUri"`
	Token        *string  `json:"token"`
}

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=postgres dbname=misskey sslmode=disable host=db")
	if err != nil {
		slog.Error("Error opening database", "err", err)
	}
}

func FindOneByUsername(username string) (*MiUser, error) {
	rows, err := db.Query(`SELECT id, "updatedAt", "lastFetchedAt", "lastActiveDate", "hideOnlineStatus",
		"username", "usernameLower", "name", "followersCount", "followingCount", "movedToUri", "movedAt",
		"alsoKnownAs", "notesCount", "avatarId", "bannerId", "avatarUrl", "bannerUrl", "avatarBlurhash",
		 "bannerBlurhash", "avatarDecorations", "tags", "isSuspended", "isLocked", "isBot", "isCat",
		 "isRoot", "isExplorable", "isHibernated", "isDeleted", "emojis", "host", "inbox", "sharedInbox",
		  "featured", "uri", "followersUri", "token" FROM public."user" WHERE "username" = $1`, username)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var u MiUser
		var avatarDecorations string
		var tags string
		var emojis string
		err = rows.Scan(&u.ID, &u.UpdatedAt, &u.LastFetchedAt, &u.LastActiveDate, &u.HideOnlineStatus, &u.Username,
			&u.UsernameLower, &u.Name, &u.FollowersCount, &u.FollowingCount, &u.MovedToURI, &u.MovedAt, &u.AlsoKnownAs,
			&u.NotesCount, &u.AvatarID, &u.BannerID, &u.AvatarURL, &u.BannerURL, &u.AvatarBlurHash, &u.BannerBlurHash,
			&avatarDecorations, &tags, &u.IsSuspended, &u.IsLocked, &u.IsBot, &u.IsCat, &u.IsRoot, &u.IsExplorable,
			&u.IsHibernated, &u.IsDeleted, &emojis, &u.Host, &u.Inbox, &u.SharedInbox, &u.Featured, &u.URI, &u.FollowersURI, &u.Token)
		if err != nil {
			return nil, err
		}

		// Parse avatarDecorations
		json.Unmarshal([]byte(avatarDecorations), &u.AvatarDecorations)

		// Parse tags
		// json.Unmarshal([]byte(tags), &u.Tags)

		// Parse emojis

		return &u, nil
	}

	return nil, nil
}
