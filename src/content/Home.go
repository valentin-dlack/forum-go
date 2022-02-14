package content

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

// servHome : //* Page d'acceuil
func ServeHome(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)
	color := RandomColor()

	db, err := sql.Open("sqlite3", "database/database.db")
	CheckErr(err)

	allPosts, _ := db.Query("SELECT * FROM Posts ORDER BY id DESC LIMIT 3")

	var post POSTINFO
	var mostRecent []POSTINFO
	var since string
	var post_id int
	var title string
	var categories string
	var body string
	var user_id int
	var image string
	var likes int
	var comments_nb int
	for allPosts.Next() {
		err = allPosts.Scan(&post_id, &title, &categories, &body, &user_id, &image, &likes, &comments_nb, &since)
		CheckErr(err)

		cat := strings.Split(categories, ";")
		var tabCategories []CATEGORIES
		for _, x := range cat {
			catephemere := CATEGORIES{
				Cat:   x,
				Color: color[x],
			}
			tabCategories = append(tabCategories, catephemere)
		}

		tabusers, err := db.Query("SELECT * FROM Users")
		if err != nil {
			fmt.Println(err.Error())
		}
		var userinfo INFO
		var userID int
		var username string
		var email string
		var since string
		var description string
		var password string
		var country string
		var mod int
		for tabusers.Next() {
			err = tabusers.Scan(&userID, &username, &email, &since, &description, &password, &image, &country, &mod)
			CheckErr(err)

			if userID == user_id {
				userinfo = GetUser(userID)
				break
			}
		}
		tabusers.Close()

		post = POSTINFO{
			ID:         post_id,
			User_ID:    user_id,
			Title:      title,
			Body:       body,
			Image:      image,
			Categories: tabCategories,
			Likes:      likes,
			Comment_Nb: comments_nb,

			Post_User_Info: userinfo,
		}
		mostRecent = append(mostRecent, post)
	}

	allPosts.Close()

	allPosts, _ = db.Query("SELECT * FROM Posts ORDER BY likes DESC LIMIT 3")
	var mostLikes []POSTINFO
	for allPosts.Next() {
		err = allPosts.Scan(&post_id, &title, &categories, &body, &user_id, &image, &likes, &comments_nb, &since)
		CheckErr(err)

		cat := strings.Split(categories, ";")
		var tabCategories []CATEGORIES
		for _, x := range cat {
			catephemere := CATEGORIES{
				Cat:   x,
				Color: color[x],
			}
			tabCategories = append(tabCategories, catephemere)
		}

		tabusers, err := db.Query("SELECT * FROM Users")
		if err != nil {
			fmt.Println(err.Error())
		}
		var userinfo INFO
		var userID int
		var username string
		var email string
		var since string
		var description string
		var password string
		var country string
		var mod int
		for tabusers.Next() {
			err = tabusers.Scan(&userID, &username, &email, &since, &description, &password, &image, &country, &mod)
			CheckErr(err)
			if userID == user_id {
				userinfo = GetUser(userID)
				break
			}
		}
		tabusers.Close()

		post = POSTINFO{
			ID:         post_id,
			User_ID:    user_id,
			Title:      title,
			Body:       body,
			Image:      image,
			Categories: tabCategories,
			Likes:      likes,
			Comment_Nb: comments_nb,

			Post_User_Info: userinfo,
		}
		mostLikes = append(mostLikes, post)
	}

	allPosts.Close()

	db.Close()

	data := ALLINFO{
		Self_User_Info: user,
		Post_Info:      POSTINFO{},

		Post_Most_Recent: mostRecent,
		Post_Most_Likes:  mostLikes,
	}

	files := []string{"template/Home.html", "template/Common.html"}

	tmp, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error: Check template", 500)
	}

	err = tmp.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error", 500)
	}

}
