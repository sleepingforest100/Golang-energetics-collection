              ENERGETIX COLLECTION - web service for collecting energy drinks + web page 
DESCRIPTION: Nowadays, there is large amount of energy drinks, and a lot of them differs greatly. There are different tastes, various cans, distinct styles, and it's impossible to remember them 
all!. So this project can help with this problem. The application allows users to collect information about different energy drinks. It gives opportunities to add energy drinks in database, 
update information about them by id, delete energy drinks by id, find a drink by id or look through them all. The target users are people who love drinking and collecting energy drinks.
The data is stored in postgres. The server starts on the port 8080.

AUTHORS: Polina Batova SE-2212, Igor Letunovskii SE-2212, Madina Aitzhanova SE-2212

SCREENSHOT
![screenshot](https://github.com/PollyBreak/Golang-energetics-collection/assets/88556120/fb781850-c834-4984-b686-1bdf3951254d)


LAUNCH INSTRUCTIONS:
  1. Dowload and install Postgres, create connection, create database NAMED "energetix"
  2. Clone this repository.
  3. Now we should use migrations for creating tables. Firstly, you need to dowload 'migrate'. For example, you can use Powershell and Scoop in Windows. For this use these instructions:
       https://www.freecodecamp.org/news/database-migration-golang-migrate/ (Only "the Setup and Installation" Chapter)
       https://scoop.sh/#/ (Only for WINDOWS)
     After installing use the next command
         migrate -path database/migration/ -database "postgresql://username:secretkey@localhost:5432/energetix?sslmode=disable" -verbose up
     -path "database/migration/" should be raplaced with the path to migrations from the repository. For example, my path is "C:/GIT/Go/database/migration/"/
     username should be raplaced with your username in postgres (by default it is 'postgres')
     secretkey should be raplaced with your password in postgres
     localhost should be raplaced with the port that you use for postgres.
     After successfull executing of the command you should see something like that:
        ...
        2024/01/07 22:42:13 Finished after 383.6799ms
        2024/01/07 22:42:13 Closing source and database
     If you see these lines, the tables were created.
  4. Open the file main.go. Install in the terminal all missing imports. For this I wrote the next commands in my VisualStudio Code Terminal.
        go get -u gorm.io/gorm 
        go get -u gorm.io/driver/postgres
        go get -u github.com/lib/pq
        go get -u github.com/gorilla/handlers
        go get -u github.com/gorilla/mux
        go get -u golang.org/x/time/rate
        go get -u github.com/sirupsen/logrus
  6. Change the password and user for db connection with yours in the function initDB(). It's very important!
  7. Launch the application, write "go run main.go" in the terminal. Agree with your FireWall if it's neccessary.
  8. Open "http://localhost:8080/index-go.html" in your browser.
  9. Scroll page a bit lower and click on "Moderation" to the right side of the page. There you can see options to create and to find by id.
  10. For updating and deleting you can click on the energetics cards.
     ![screenshot2](https://github.com/PollyBreak/Golang-energetics-collection/assets/88556120/6c63c443-c7aa-4427-a9f5-264b5235ecdf)



       
TOOLS STACK: 
  BACKEND: Golang, Postgres, GORM (ORM library), golang-migrate, logrus, /x/time/rate
  FRONTEND:HTML, CSS, JavaScript

LINKS TO SOURCES
  Rwitesh Bera. How to Perform Database Migrations using Go Migrate. URL: https://www.freecodecamp.org/news/database-migration-golang-migrate/ 
  GORM Documentation. URL: https://gorm.io/docs/index.html
  One-to-one in GO-language GORM. URL: https://programmerall.com/article/23262004501/#google_vignette
  Hugo Johnsson. REST-API with Golang and Mux. URL: https://hugo-johnsson.medium.com/rest-api-with-golang-and-mux-e934f581b8b5


