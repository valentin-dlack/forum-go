package content

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)


// Delete : "Supprime le UUID du compte qui se déconnect"
func Delete(w http.ResponseWriter, r *http.Request, id int) {
	db, err := sql.Open("sqlite3", "database/database.db")
	CheckErr(err)

	stmt, err := db.Prepare("delete from sessions where user_id=?")
	CheckErr(err)

	res, err := stmt.Exec(id)
	CheckErr(err)

	_, err = res.RowsAffected()
	CheckErr(err)

	db.Close()

	c := http.Cookie{Name: "sessionLog", Value: "", MaxAge: -1}
	http.SetCookie(w, &c)
	http.Redirect(w, r, "/login", 301)
}

func CookieCreation(w http.ResponseWriter, id int) {
	db, err := sql.Open("sqlite3", "database/database.db")
	CheckErr(err)
	u1 := uuid.NewV4()
	fmt.Printf("UUID : %s ; User_id : %d\n", u1, id)
	datab, err := db.Prepare("INSERT INTO sessions (user_id, uuid) VALUES (?, ?)")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server Error", 500)
	}
	_, err = datab.Exec(id, u1)
	if err != nil {
		fmt.Println(err)
	}

	db.Close()
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "sessionLog", Value: String(u1), Expires: expiration}
	http.SetCookie(w, &cookie)
}

func CheckSession(r *http.Request) (bool, int) {
	db, err := sql.Open("sqlite3", "database/database.db")
	CheckErr(err)
	cookie, _ := r.Cookie("sessionLog")
	ok := false
	var id int
	var uuid string
	if cookie != nil {
		dataCookie, err := db.Query("SELECT * FROM sessions")
		CheckErr(err)

		for dataCookie.Next() {
			err = dataCookie.Scan(&id, &uuid)
			CheckErr(err)
			if cookie.Value == uuid {
				ok = true
				break
			}
		}
		dataCookie.Close()
	}

	db.Close()
	return ok, id
}