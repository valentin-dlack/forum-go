package content

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func CreationPost(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)

	if user.UserName != "" {
		r.ParseMultipartForm(10 << 20) //max size 10Mb (5mb for the pf)
		var Post POSTINFO
		var tabCat []CATEGORIES

		color := RandomColor()

		//sport, anime/manga, economie, jeux vidéo, informatique, voyages, NEW, paranormal.

		db, err := sql.Open("sqlite3", "database/database.db")
		CheckErr(err)

		if r.Method == "POST" {
			datab, err := db.Prepare("INSERT INTO Posts (title, categories, body, user_id, image, likes, comment_nb, since) VALUES (?,?,?,?,?,?,?, ?)")
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Server Error", 500)
			}
			title := r.FormValue("title")
			body := r.FormValue("body")
			var image string
			file, handler, err := r.FormFile("myFile")
			if err != nil {
				image = r.FormValue("myFile")
				fmt.Println(err)
			} else {
				// defer file.Close()
				// fmt.Printf("Uploaded File: %+v\n", strings.ReplaceAll(handler.Filename, " ", "-"))
				// fmt.Printf("File Size: %+v\n", handler.Size)
				// fmt.Printf("MIME Header: %+v\n", handler.Header)

				absPath, _ := filepath.Abs("../src/assets/posts/" + strings.ReplaceAll(handler.Filename, " ", "-"))

				resFile, err := os.Create(absPath)
				if err != nil {
					fmt.Print(w, err)
				}
				defer resFile.Close()

				io.Copy(resFile, file)
				defer resFile.Close()
				fmt.Print("File uploaded")

				image = "../assets/posts/" + strings.ReplaceAll(handler.Filename, " ", "-")
			}

			likes := 0
			comment_nb := 0
			loc, _ := time.LoadLocation("Europe/Paris")
			pretime := time.Now().In(loc)
			since := pretime.String()[:19]
			var categoriesCheck string
			for categorie := range color {
				if r.FormValue(categorie) != "" {
					categoriesCheck += categorie + ";"
				}
			}

			if categoriesCheck == "" {
				categoriesCheck += "autre;"
			}

			if title != "" && body != "" && categoriesCheck != "" {
				user_id := user.ID
				var tabAllCat []CATEGORIES

				_, err := datab.Exec(title, categoriesCheck, body, user_id, image, likes, comment_nb, since)
				if err != nil {
					fmt.Println(err.Error())
				}

				tabCategoriesCheck := strings.Split(categoriesCheck, ";")
				for _, x := range tabCategoriesCheck {
					oneCategorie := CATEGORIES{
						Cat:   x,
						Color: color[x],
					}
					tabCat = append(tabCat, oneCategorie)
				}
				for x := range color {
					oneCategorie := CATEGORIES{
						Cat:   x,
						Color: color[x],
					}
					tabAllCat = append(tabAllCat, oneCategorie)
				}
				Post = POSTINFO{
					User_ID:       user.ID,
					Title:         title,
					Body:          body,
					Image:         image,
					Categories:    tabCat,
					Since:         since,
					AllCategories: tabAllCat,
				}

				uId := strconv.Itoa(user.ID)

				newPost, err := db.Query("SELECT * FROM Posts WHERE user_id=" + uId + " ORDER BY id DESC LIMIT 1")
				if err != nil {
					fmt.Println(err.Error())
				}
				var id string
				var categories string
				for newPost.Next() {
					err = newPost.Scan(&id, &title, &categories, &body, &user_id, &image, &likes, &comment_nb, &since)
					CheckErr(err)
				}
				newPost.Close()

				http.Redirect(w, r, "/post?id="+id, 301)
			}
		} else {
			for x := range color {
				oneCategorie := CATEGORIES{
					Cat:   x,
					Color: color[x],
				}
				tabCat = append(tabCat, oneCategorie)
			}
			Post = POSTINFO{
				Categories: tabCat,
			}
		}
		data := ALLINFO{
			Self_User_Info: user,
			Post_Info:      Post,
		}

		files := []string{"template/CreatePost.html", "template/Common.html"}
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
	} else {
		http.Redirect(w, r, "/login", 301)
	}
}
