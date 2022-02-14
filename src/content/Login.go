package content

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

// Login : Permet de se connecter à un compte
func Login(w http.ResponseWriter, r *http.Request) {
	userinfo := INFO{
		ID:          0,
		Email:       "",
		PassWord:    "",
		UserName:    "",
		Since:       "",
		Description: "",
		Image:       "",
		Country:     "",
		Login:       false,
		Msg:         "Vous êtes déconnecté",
	}

	db, err := sql.Open("sqlite3", "database/database.db")
	CheckErr(err)

	cExist, id := CheckSession(r)

	if r.Method == "POST" {
		if cExist {
			Delete(w, r, id)
		} else {
			test, err := db.Query("SELECT * FROM Users")
			if err != nil {
				fmt.Println(err.Error())
			}
			mailfound := false
			var id int
			var email string
			var Password string
			var username string
			var since string
			var description string
			var image string
			var country string
			var modB int
			var mod bool
			for test.Next() {
				err = test.Scan(&id, &username, &email, &since, &description, &Password, &image, &country, &modB)
				CheckErr(err)
				if email == r.FormValue("mail") {
					mailfound = true
					break
				}
			}
			mdp := r.FormValue("password")

			test.Close()
			mod = IntToBoolAdmin(modB)
			if mailfound {
				fmt.Print("Into mailfound")
				cryptedPassword := []byte(Password)
				if bcrypt.CompareHashAndPassword(cryptedPassword, []byte(mdp)) == nil {
					CookieCreation(w, id)
					userinfo = INFO{
						ID:          id,
						Email:       email,
						PassWord:    mdp,
						UserName:    username,
						Since:       since,
						Description: description,
						Image:       image,
						Country:     country,
						Admin:       mod,
						Login:       true,
						Msg:         "Vous êtes connecté en tant que " + username,
					}

					http.Redirect(w, r, "/", 301)
				} else {
					userinfo = INFO{
						Msg: "Le mot de passe est invalide",
					}
				}
			} else {
				userinfo = INFO{
					Msg: "Ce mail n'est pas enregistré: veuillez vous inscrire",
				}
			}
		}

	} else {
		if cExist {
			userinfo = GetSession(r)
		}
	}

	data := ALLINFO{
		Self_User_Info: userinfo,
	}
	files := []string{"template/Connexion.html", "template/Common.html"}
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
	db.Close()
}
