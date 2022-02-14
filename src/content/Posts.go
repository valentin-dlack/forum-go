package content

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"text/template"
)

// AllPosts : kzdnzndz
func AllPosts(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)
	color := RandomColor()
	var tabCat []CATEGORIES
	var tabATrier []string
	for categorie := range color {
		tabATrier = append(tabATrier, categorie)
	}
	sort.Strings(tabATrier)
	for _, categorie := range tabATrier {
		oneCategorie := CATEGORIES{
			Cat:   categorie,
			Color: color[categorie],
		}
		tabCat = append(tabCat, oneCategorie)
	}

	var all_Post []POSTINFO
	var post_info POSTINFO

	db, err := sql.Open("sqlite3", "database/database.db")

	categorieCheck := ""
	var result ALLINFO
	var searching = false
	if r.Method == "POST" {
		if r.FormValue("search") != ""{
			result = SearchData(r.FormValue("search"))
			searching = true
		}else{
			for categorie := range color {
				if r.FormValue(categorie) != "" {
					categorieCheck = categorie
				}
			}
		}


	}
	var data ALLINFO
	if searching{
		postInfo := POSTINFO{
			AllCategories: tabCat,
		}
	
		data = ALLINFO{
			Self_User_Info: user,
			Post_Info:      postInfo,
	
			All_User:  result.All_User,
			All_Posts: result.All_Posts,
		}
	}else{
		post, err := db.Query("SELECT * FROM Posts ORDER BY id DESC")

		var since string
		var id int
		var user_id int
		var title string
		var body string
		var image string
		var likes int
		var comment_nb int
		var categories string
		var userinfo INFO
		for post.Next() {
			err = post.Scan(&id, &title, &categories, &body, &user_id, &image, &likes, &comment_nb, &since)
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
			catCheck := false
			for _, y := range tabCategories {
				if y.Cat == categorieCheck {
					catCheck = true
					continue
				}
			}
			if catCheck == true {
				userinfo = GetUser(user_id)

				post_info = POSTINFO{
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
				all_Post = append(all_Post, post_info)
			}

		}
		post.Close()

		db.Close()

		postInfo := POSTINFO{
			AllCategories: tabCat,
		}

		data = ALLINFO{
			Self_User_Info: user,
			Post_Info:      postInfo,

			All_User:  result.All_User,
			All_Posts: all_Post,
		}
	}

	files := []string{"template/Posts.html", "template/Common.html"}

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
