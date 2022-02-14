package content

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Println("() ====================== ( ! ERROR ! ) ====================== ()")
		panic(err)
	}
}

func String(u uuid.UUID) string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

func RandomColor() map[string]string {
	allColor := map[string]string{
		"anime/manga":  "#D50C2E",
		"autre":        "#ccd1d1",
		"culture":      "#E15256",
		"economie":     "#C3C020",
		"informatique": "#19A9D1",
		"jeux vidéos":  "#23C009",
		"NEWS":         "#ff5733",
		"paranormal":   "#070709",
		"sport":        "#9D84C9",
		"voyage":       "#00FF12",
	}
	return allColor
}


func SearchData(search string) ALLINFO{
	db, err := sql.Open("sqlite3", "database/database.db")
	CheckErr(err)

	var allUsers []INFO

	users, err := db.Query("SELECT * FROM Users ORDER BY id DESC")

	var currentlyUser INFO
	var email string
	var password string
	var username string
	var description string
	var country string
	var mod int
	var id int
	var image string
	var since string
	for users.Next() {
		err = users.Scan(&id, &username, &email, &since, &description, &password, &image, &country, &mod)
		CheckErr(err)
		if strings.Contains(strings.ToLower(username), strings.ToLower(search)){
			currentlyUser = GetUser(id)
			allUsers = append(allUsers, currentlyUser)
		}
	}
	users.Close()



	posts, _ := db.Query("SELECT * FROM Posts ORDER BY id DESC")

	color := RandomColor()
	var allPost []POSTINFO
	var user_id int
	var title string
	var body string
	var likes int
	var comment_nb int
	var categories string
	var userinfo INFO
	for posts.Next() {
		err = posts.Scan(&id, &title, &categories, &body, &user_id, &image, &likes, &comment_nb, &since)
		CheckErr(err)
		cat := strings.Split(categories, ";")
		var tabCategories []CATEGORIES
		for _, onecategorie := range cat {
			catephemere := CATEGORIES{
				Cat:   onecategorie,
				Color: color[onecategorie],
			}
			tabCategories = append(tabCategories, catephemere)
		}
			userinfo = GetUser(user_id)

			post_info := POSTINFO{
				ID:         id,
				User_ID:    user_id,
				Title:      title,
				Body:       body,
				Image:      image,
				Categories: tabCategories,
				Likes:      likes,
				Comment_Nb: comment_nb,

				Post_User_Info: userinfo,
			}
			if strings.Contains(strings.ToLower(userinfo.UserName), strings.ToLower(search)) || strings.Contains(strings.ToLower(title), strings.ToLower(search)) || strings.Contains(strings.ToLower(body), strings.ToLower(search))|| strings.Contains(strings.ToLower(categories), strings.ToLower(search)){
				allPost = append(allPost, post_info)
			}

	}
	posts.Close()



	db.Close()

	result := ALLINFO{
		All_User: allUsers,
		All_Posts: allPost,
	}

	return result
}