package content

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)


func GetSession(r *http.Request) INFO {
	var userinfo INFO
	fmt.Println("Get Session")
	cExist, idSession := CheckSession(r)
	if cExist {

		db, err := sql.Open("sqlite3", "database/database.db")
		CheckErr(err)

		tabusers, err := db.Query("SELECT * FROM Users")
		if err != nil {
			fmt.Println(err.Error())
		}
		var id int
		var email string
		var password string
		var username string
		var since string
		var description string
		var image string
		var country string
		var mod int

		for tabusers.Next() {
			err = tabusers.Scan(&id, &username, &email, &since, &description, &password, &image, &country, &mod)
			CheckErr(err)

			if id == idSession {

				user := GetUser(id)
				posts := GetPost(user)
				userinfo = INFO{
					ID:          id,
					Email:       user.Email,
					PassWord:    user.PassWord,
					UserName:    user.UserName,
					Since:       user.Since,
					Description: user.Description,
					Image:       user.Image,
					Country:     user.Country,
					Admin:       user.Admin,
					Modo:        user.Modo,
					Login:       true,
					AllPosts:    posts,
				}
				break
			}
		}
		tabusers.Close()
		db.Close()
	}
	return userinfo
}

func GetUser(id int) INFO {
	db, err := sql.Open("sqlite3", "database/database.db")
	if err != nil {
		fmt.Print(err)
	}
	tabusers, err := db.Query("SELECT * FROM Users where id=" + strconv.Itoa(id))
	if err != nil {
		fmt.Println(err.Error())
	}
	var userinfo INFO
	var userAllPost []POSTINFO
	var userID int
	var username string
	var image string
	var email string
	var description string
	var password string
	var country string
	var since string
	var mod int
	for tabusers.Next() {
		err = tabusers.Scan(&userID, &username, &email, &since, &description, &password, &image, &country, &mod)
		CheckErr(err)
		if userID == id {
			userinfo = INFO{
				ID:          id,
				Email:       email,
				PassWord:    password,
				UserName:    username,
				Since:       since,
				Description: description,
				Image:       image,
				Country:     country,
				Mod:         mod,
			}
			userAllPost = GetPost(userinfo)

			break
		}
	}

	admin := IntToBoolAdmin(userinfo.Mod)
	modo := IntToBoolModo(userinfo.Mod)

	userinfo = INFO{
		ID:          id,
		Email:       email,
		PassWord:    password,
		UserName:    username,
		Since:       since,
		Description: description,
		Image:       image,
		Country:     country,
		Admin:       admin,
		Modo:        modo,
		AllPosts:    userAllPost,
	}
	tabusers.Close()
	db.Close()
	return userinfo
}

func IntToBoolAdmin(mod int) bool {
	if mod == 2 {
		return true
	} else {
		return false
	}
}

func IntToBoolModo(mod int) bool {
	if mod == 1 {
		return true
	} else {
		return false
	}
}

func GetPost(user INFO) []POSTINFO {
	var all_Post []POSTINFO
	db, err := sql.Open("sqlite3", "database/database.db")

	post, err := db.Query("SELECT * FROM Posts WHERE user_id=" + strconv.Itoa(user.ID) + " ORDER BY id DESC")

	color := RandomColor()
	var since string
	var user_id int
	var id string
	var title string
	var body string
	var image string
	var likes int
	var comment_nb int
	var categories string
	for post.Next() {
		err = post.Scan(&id, &title, &categories, &body, &user_id, &image, &likes, &comment_nb, &since)
		CheckErr(err)
		idInt, _ := strconv.Atoi(id)
		cat := strings.Split(categories, ";")
		var tabCategories []CATEGORIES
		for _, x := range cat {
			catephemere := CATEGORIES{
				Cat:   x,
				Color: color[x],
			}
			tabCategories = append(tabCategories, catephemere)
		}
		post_info := POSTINFO{
			ID:             idInt,
			User_ID:        user_id,
			Title:          title,
			Body:           body,
			Image:          image,
			Categories:     tabCategories,
			Likes:          likes,
			Comment_Nb:     comment_nb,
			Since:          since,
			Post_User_Info: user,
		}
		all_Post = append(all_Post, post_info)
	}
	post.Close()

	db.Close()
	return all_Post
}