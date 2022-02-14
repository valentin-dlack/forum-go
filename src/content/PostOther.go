package content

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)


func EditPost(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)
	color := RandomColor()

	postID := r.FormValue("id")
	if user.UserName != "" {
		var data ALLINFO
		var post_id int
		var title string
		var categories string
		var body string
		var user_id int
		var image string
		var likes int
		var comments_nb int
		var since string

		db, err := sql.Open("sqlite3", "database/database.db")
		CheckErr(err)

		Post, err := db.Query("SELECT * FROM Posts WHERE id=" + postID)
		if err != nil {
			fmt.Println(err.Error())
		}
		for Post.Next() {
			err = Post.Scan(&post_id, &title, &categories, &body, &user_id, &image, &likes, &comments_nb, &since)
			CheckErr(err)
		}
		Post.Close()

		tabCategories := strings.Split(categories, ";")
		var tabCat []CATEGORIES
		for _, x := range tabCategories {
			if x != "" {
				oneCategorie := CATEGORIES{
					Cat:   x,
					Color: color[x],
				}
				tabCat = append(tabCat, oneCategorie)
			}
		}

		var tabAllCat []CATEGORIES
		for x, _ := range color {
			var check string
			for _, y := range tabCategories {
				if y == x {
					check = "checked"
					break
				}
			}
			oneCategorie := CATEGORIES{
				Cat:   x,
				Color: color[x],
				Check: check,
			}
			tabAllCat = append(tabAllCat, oneCategorie)
		}

		post_info := POSTINFO{
			ID:            post_id,
			User_ID:       user_id,
			Title:         title,
			Body:          body,
			Image:         image,
			Categories:    tabCat,
			AllCategories: tabAllCat,
			Likes:         likes,
			Comment_Nb:    comments_nb,
		}

		if user.ID == post_info.User_ID {
			var tabCat []CATEGORIES

			if r.Method == "POST" {

				//Récupération des nouvelles entrés
				newTitle := r.FormValue("title")
				newBody := r.FormValue("body")
				newImage := r.FormValue("Image")
				var newCategories string
				for x := range color {
					if r.FormValue(x) != "" {
						newCategories += x + ";"
					}
				}
				tabCategoriesCheck := strings.Split(newCategories, ";")
				for _, x := range tabCategoriesCheck {
					oneCategorie := CATEGORIES{
						Cat:   x,
						Color: color[x],
					}
					tabCat = append(tabCat, oneCategorie)
				}
				var Title string
				if newTitle != "" {
					Title = newTitle
				} else {
					Title = title
				}
				var Body string
				if body != "" {
					Body = newBody
				} else {
					Body = body
				}
				var Categories string
				if newCategories != "" {
					Categories = newCategories
				} else {
					Categories = categories
				}
				var Image string
				if newImage != "" {
					Image = newImage
				} else {
					Image = image
				}

				edit, _ := db.Prepare("UPDATE Posts SET title=?, categories=?, body=?, image=? WHERE id=" + postID)

				_, err := edit.Exec(Title, Categories, Body, Image)
				if err != nil {
					fmt.Println(err.Error())
				}

				post_info = POSTINFO{
					ID:         post_id,
					User_ID:    user_id,
					Title:      Title,
					Body:       Body,
					Image:      Image,
					Categories: tabCat,
				}

				edit.Close()
				http.Redirect(w, r, "/post?id="+strconv.Itoa(post_info.ID), 301)
			}
			data := ALLINFO{
				Self_User_Info: user,
				Post_Info:      post_info,
			}

			files := []string{"template/EditPost.html", "template/Common.html"}
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
		}else{
			files := []string{"template/404.html"}
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
		db.Close()
	} else {
		http.Redirect(w, r, "/login", 301)
	}
}

func DeletePost(id string, user INFO) {

	postID := id
	if user.UserName != "" {
		var post_id int
		var title string
		var categories string
		var body string
		var user_id int
		var image string
		var likes int
		var comments_nb int
		var since string

		db, err := sql.Open("sqlite3", "database/database.db")
		CheckErr(err)
		post, err := db.Query("SELECT * FROM Posts WHERE id=" + postID)
		if err != nil {
			fmt.Println(err.Error())
		}

		CheckErr(err)
		for post.Next() {
			err = post.Scan(&post_id, &title, &categories, &body, &user_id, &image, &likes, &comments_nb, &since)
			CheckErr(err)
		}
		post.Close()

		post_info := POSTINFO{
			ID:         post_id,
			User_ID:    user_id,
			Title:      title,
			Body:       body,
			Image:      image,
			Likes:      likes,
			Comment_Nb: comments_nb,
		}

		if user.ID == post_info.User_ID || user.Admin || user.Modo{

			del, _ := db.Prepare("DELETE from Posts WHERE id=?")

			res, err := del.Exec(post_info.ID)
			CheckErr(err)

			_, err = res.RowsAffected()
			CheckErr(err)

			del.Close()

			comment, err := db.Prepare("DELETE from Comments WHERE post_id=?")

			CheckErr(err)

			res, err = comment.Exec(post_info.ID)
			CheckErr(err)

			_, err = res.RowsAffected()
			CheckErr(err)

			comment.Close()

			Like, err := db.Prepare("DELETE from Likes WHERE post_id=?")

			CheckErr(err)

			res, err = Like.Exec(post_info.ID)
			CheckErr(err)

			_, err = res.RowsAffected()
			CheckErr(err)

			Like.Close()
		}
		db.Close()
	}
}

func DeleteCommentaire(id string,post_id string,  user INFO) {

	
	if user.UserName != "" {
		// postID := post_id
		// var body string

		db, err := sql.Open("sqlite3", "database/database.db")
		CheckErr(err)
		// post, err := db.Query("SELECT * FROM Comments WHERE id=" + id)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }

		// CheckErr(err)
		// for post.Next() {
		// 	err = post.Scan(&id,&body, &user_id,&post_id, &since)
		// 	CheckErr(err)
		// }
		// post.Close()

		// comment := COMMENT{
		// 	ID:         id,
		// 	User_ID:    user_id,
		// 	User_Info: user,
		// 	Post_ID: post_id,
		// 	Body:       body,
		// }

		if /*user.ID == post_info.User_ID ||*/ user.Admin || user.Modo{

			del, _ := db.Prepare("DELETE from Comments WHERE id=?")

			res, err := del.Exec(id)
			CheckErr(err)

			_, err = res.RowsAffected()
			CheckErr(err)

			del.Close()

		}
		db.Close()
	}
}