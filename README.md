# Ytrack Forum Project

> Dautrement Valentin, Arthur ABADIE, Matteo FERREIRA, Costa REYPES

## Project Presentation

The "Forum" Project is a a simple complete website project, based on a forum like reddit or else. The goal is to make a forum easy to use, aesthetic but also secure with all the basic features that a forum should have: Posts (with images), comments, likes and even a moderation system !

## Setup
> You need to get the source code with a `git clone [git link]` before running this setup.

* To use the program, go into /forum/src and you'll need to type `go run server.go` on a terminal or cmd.

![Powershell console with "go run ./server.go"](https://i.imgur.com/0BSQNku.png)

* Now you need to authorize the firewall

![Athorize the firewall](https://i.imgur.com/VQZIko5.png)

*  You can finally go to your favorite web browser and type the address indicated in the terminal *(here `localhost:4444`)* in the search bar 

![](https://i.imgur.com/Hfx8nbT.png)

---

## Usage 

---

### Login / Register :

If you have an account, you can login, otherwise you can always create a new account by clicking on "IJe n'ai pas de compte".

![](https://i.imgur.com/9MbV9mk.png)

![](https://i.imgur.com/JZh516E.png)

---

### Create Post :

Now that you're logged in, you have a bunch of new options, you can, for example, create a new post. You just have to click on the "+" button at the bottom of the page !

![](https://i.imgur.com/rwvK1tI.png)

And now you are on the create post page, you can create your post with choosen categories, images, and others !

![](https://i.imgur.com/wi6oGQT.png)

---

### Post Options :

Now that your post is online, people can interact with it, comment it, like it or report it, and YOU can also delete your post. 

![](https://i.imgur.com/cvt7oJu.png)

---

### Profile Options :

Since you are connected, you can modify some things in your profile like the description, the name, the profile picture, ect...  

You just need to click on "Edit Profile"  

![](https://i.imgur.com/oMuo2ss.png) ![](https://i.imgur.com/eI1PmLu.png)

---

### Moderation :

The final section is for moderation, moderators can delete posts and base users, and administrators can upgrade users to mods, deletes posts, users, mods ect...  
For those we have specials control panels like this :  

**Profile control panel :**  

![](https://i.imgur.com/weZFgQe.png)

You can access to the profile of who you want and then you have options that appear (Delete or Promote)

![](https://i.imgur.com/wGgmhkJ.png)

**Post Control Panel :**  

Here you can directly delete the post from the panel, but there's also a button to go to the post for more informations, there is also an image preview (if the post have an image)

![](https://i.imgur.com/o8AteeB.png)


---

## Credits
Thanks for checking our project ! 

- Dautrement Valentin : Backend (Golang), Front (html, bootstrap and JS)
- Arthur ABADIE : Front (HTML, CSS, Bootstrap and a bit of JS)
- Matteo FERREIRA : Front (Base Design, HTML, CSS, Bootstrap), Databases (SQLite)
- Costa REYPES : Backend (Golang), Front to Back connection (Go Templates).  