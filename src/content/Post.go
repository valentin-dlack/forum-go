package content

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//OnePost : Page pour un seul post
func OnePost(w http.ResponseWriter, r *http.Request) {
	userInfo := GetSession(r)
	color := RandomColor()

	likeNow := ""
	likes := 0
	commentnb := 0

	post_id := r.FormValue("id")
	upost_id, err := strconv.Atoi(post_id)

	db, _ := sql.Open("sqlite3", "database/database.db")
	like, err := db.Query("SELECT * FROM Likes WHERE user_id=" + strconv.Itoa(userInfo.ID))
	if err != nil {
		fmt.Println(err.Error())
	}

	var id int
	var idPost int
	var idUser int
	var since string
	var dejaLike bool = false
	for like.Next() {
		err = like.Scan(&id, &idPost, &idUser, &since)
		CheckErr(err)
		if idPost == upost_id {
			dejaLike = true
			break
		}
	}

	if dejaLike {
		likeNow = "checked"
	}

	like.Close()

	db.Close()

	

	//Récupération du nouveau commentaire
	if r.Method == "POST" {
		db, _ := sql.Open("sqlite3", "database/database.db")
		comment := r.FormValue("comment")

		changement := false
		if (dejaLike && r.FormValue("Liker")=="Liker") || (!dejaLike && r.FormValue("Liker") == "") {
			changement=false
		}else{
			changement=true
		}
		if userInfo.UserName != ""{
			if comment != ""  && !changement {
				datab, err := db.Prepare("INSERT INTO Comments (body, user_id,post_id,since) VALUES (?,?,?,?)")
				if err != nil {
					fmt.Println(err)
					http.Error(w, "Server Error", 500)
				}

				user_id := userInfo.ID
				post_id := upost_id
				loc, _ := time.LoadLocation("Europe/Paris")
				pretime := time.Now().In(loc)
				since := pretime.String()[:19]
				_, err = datab.Exec(comment, user_id, post_id, since)
				if err != nil {
					fmt.Println(err)
				}
				datab.Close()
			} else if r.FormValue("deleteButton") != "" {
				DeletePost(post_id, userInfo)
				http.Redirect(w,r,"/posts", 301)
			} else if changement{
				if !dejaLike {

					loc, _ := time.LoadLocation("Europe/Paris")
					pretime := time.Now().In(loc)
					since := pretime.String()[:19]
					datab, err := db.Prepare("INSERT INTO Likes (user_id,post_id,since) VALUES (?,?,?)")
					if err != nil {
						fmt.Println(err)
						http.Error(w, "Server Error", 500)
					}
					user_id := userInfo.ID
					post_id := upost_id
					_, err = datab.Exec(user_id, post_id, since)
					if err != nil {
						fmt.Println(err)
					}
					datab.Close()

					dataPost, err := db.Query("SELECT * FROM Posts WHERE id=" + strconv.Itoa(post_id))
					if err != nil {
						fmt.Println(err.Error())
					}
					likes = 0
					var title string
					var categories string
					var body string
					var image string
					var comments_nb int
					for dataPost.Next() {
						err = dataPost.Scan(&post_id, &title, &categories, &body, &user_id, &image, &likes, &comments_nb, &since)
						CheckErr(err)
					}
					dataPost.Close()

					likeNow = "checked"

				} else {
					upost_id, err := strconv.Atoi(post_id)
					CheckErr(err)
					stmt, err := db.Prepare("delete from Likes where user_id=? AND post_id=?")
					CheckErr(err)

					res, err := stmt.Exec(userInfo.ID, upost_id)
					CheckErr(err)

					_, err = res.RowsAffected()
					CheckErr(err)

					stmt.Close()

					likeNow = ""
				}
			}else if r.FormValue("commentDeleteButton") !=""{
				com_id := r.FormValue("commentDeleteButton")
				DeleteCommentaire(com_id, post_id, userInfo)
			}
		} else {
			http.Redirect(w, r, "/login", 301)
		}
		db.Close()
	}

	db, _ = sql.Open("sqlite3", "database/database.db")
	//récupération de tout les commentaires liés au post
	var title string
	var body string
	var image string
	var categories string
	var comments_nb int
	var allComments []COMMENT
	var user_id int
	var deletable bool

	getComment, err := db.Query("SELECT * FROM Comments WHERE post_id=" + post_id)
	if err != nil {
		fmt.Println(err.Error())
	}
	var id_Post int
	var comment_id int
	var bodyComment string
	for getComment.Next() {
		err = getComment.Scan(&comment_id, &bodyComment, &user_id, &id_Post, &since)
		CheckErr(err)
		if upost_id == id_Post {
			user_comment := GetUser(user_id)

			oneComment := COMMENT{
				ID:        comment_id,
				User_ID:   user_id,
				User_Info: user_comment,
				Post_ID:   id_Post,
				Body:      bodyComment,
			}

			allComments = append(allComments, oneComment)

		}
	}
	getComment.Close()

	likes = 0
	dataLikes, _ := db.Query("SELECT * FROM Likes WHERE post_id=" + post_id)
	for dataLikes.Next() {
		err = dataLikes.Scan(&id, &upost_id, &user_id, &since)
		CheckErr(err)
		likes++
	}
	dataLikes.Close()

	datab, err := db.Prepare("UPDATE Posts SET likes=? WHERE id=" + post_id)
	CheckErr(err)
	datab.Exec(likes)
	datab.Close()

	commentnb = 0
	dataComment, _ := db.Query("SELECT * FROM Comments WHERE post_id=" + post_id)
	for dataComment.Next() {
		err = dataComment.Scan(&id, &body, &user_id, &upost_id, &since)
		CheckErr(err)
		commentnb++
	}
	dataComment.Close()

	datab, err = db.Prepare("UPDATE Posts SET comment_nb=? WHERE id=" + post_id)
	CheckErr(err)
	datab.Exec(commentnb)
	datab.Close()

	//récupération du post
	test, err := db.Query("SELECT * FROM Posts WHERE id=" + post_id)
	if err != nil {
		fmt.Println(err.Error())
	}
	for test.Next() {
		err = test.Scan(&post_id, &title, &categories, &body, &user_id, &image, &likes, &comments_nb, &since)
		CheckErr(err)
	}
	test.Close()

	tabCategories := strings.Split(categories, ";")
	var tabCat []CATEGORIES
	for _, categorie := range tabCategories {
		oneCategorie := CATEGORIES{
			Cat:   categorie,
			Color: color[categorie],
		}
		tabCat = append(tabCat, oneCategorie)
	}

	//Recupération des user_info du user qui a posté
	post_user_info := GetUser(user_id)
	if post_user_info.ID == userInfo.ID || userInfo.Admin {
		deletable = true
	}
	post_info := POSTINFO{
		ID:             upost_id,
		User_ID:        post_user_info.ID,
		Title:          title,
		Body:           body,
		Image:          image,
		Categories:     tabCat,
		Likes:          likes,
		Comment_Nb:     comments_nb,
		All_Comments:   allComments,
		Post_User_Info: post_user_info,
		Deletable:      deletable,
	}

	data := ALLINFO{
		Self_User_Info:      userInfo,
		Post_Info:           post_info,
		Currently_Post_Like: likeNow,
	}
	defer db.Close()

	var files []string
	if data.Post_Info.Title != "" {
		files = []string{"template/Post.html", "template/Common.html"}
	} else {
		files = []string{"template/404.html"}
	}

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
