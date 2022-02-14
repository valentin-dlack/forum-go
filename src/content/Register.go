package content

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//Register :  Permet de se créer à un compte
func Register(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "database/database.db")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error", 500)
	}
	msg := " "
	allcountry := getPays()
	if r.Method == "POST" {
		datab, err := db.Prepare("INSERT INTO Users (username, email, since, description, password, image, country, mod) VALUES (?, ?, ?, ?, ?, ?, ?,?)")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Server Error", 500)
		}
		username := r.FormValue("username")
		email := r.FormValue("mail")
		loc, _ := time.LoadLocation("Europe/Paris")
		pretime := time.Now().In(loc)
		since := pretime.String()[:19]
		description := "Pas de description..."
		password := r.FormValue("password")
		image := "https://i.imgur.com/pMtf7R9.png"
		country := r.FormValue("country")
		mod := 0
		confirm := r.FormValue("psw-confirmation")
		Crypted := []byte(password)
		Crypted, _ = bcrypt.GenerateFromPassword(Crypted, 10)

		if username != "" || email != "" || password != "" {
			if password != confirm {
				msg = "Les deux mots de passe ne sont pas identiques"
			} else {
				_, err := datab.Exec(username, email, since, description, Crypted, image, country, mod)
				if err != nil {
					if err.Error() == "UNIQUE constraint failed: Users.email" {
						msg = "Cet E-Mail est déjà utilisé par un autre utilisateur"
					} else if err.Error() == "UNIQUE constraint failed: Users.username" {
						msg = "Ce nom est déjà utilisé par un autre utilisateur"
					} else {
					fmt.Println(err.Error())
					}
				}
				http.Redirect(w, r, "/login", 301)
			}
		}
	}
	Info := INFO{
		Msg: msg,
	}

	data := ALLINFO{
		Self_User_Info: Info,
		All_Country:    allcountry,
	}

	files := []string{"template/Register.html", "template/Common.html"}
	tmp, err := template.ParseFiles(files...) //err ne sert à rien!
	err = tmp.Execute(w, data)
	CheckErr(err)

	db.Close()
}
