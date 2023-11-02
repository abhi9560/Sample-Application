package main

import (
  "SampleApplicationGo/api"
  "SampleApplicationGo/handler"
  "context"
  "database/sql"
  //  "encoding/json"
  "errors"
  "image/png"
  //"strconv"
  "path/filepath"
  "fmt"
  "bytes"
  "github.com/foolin/echo-template"
  "github.com/go-pg/pg"
  _ "github.com/go-sql-driver/mysql"
  "github.com/gocql/gocql"
  "github.com/gorilla/mux" 
  "log"
  "github.com/boombuler/barcode"
  "github.com/boombuler/barcode/qr"
  "github.com/joho/godotenv"
  "github.com/labstack/echo"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "github.com/go-echarts/statsview"
// "strings"
  "html/template"
  "io"
  //"sync"
  "os/signal"
  "syscall"
  "net/http"
  _"net/http/pprof"
  "os"
  _"plugin"
)

// Define the template registry struct
type TemplateRegistry struct {
  templates map[string]*template.Template
}


// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
  tmpl, ok := t.templates[name]
  if !ok {
    err := errors.New("Template not found -> " + name)
    return err
  }
  return tmpl.ExecuteTemplate(w, "base.html", data)
}

type Pet struct {
  name    string
  owner   string
  species string
  sex     string
  birth   string
}
type Trainer struct {
  Name string
  Age  int
  City string
}

var (
  db *sql.DB
)

type Student struct {
  Rollno int64
  Name   string
}

func (s Student) String() string {
  return fmt.Sprintf("Student<%d %s>", s.Rollno, s.Name)
}

func pingAndQueryDB(db *sql.DB, ctx context.Context) {
  
    status := "Connection is up"
    if err := db.PingContext(ctx); err != nil {
      status = "Connection is down"
    }
    log.Println(status)
  

  rows, err := db.QueryContext(ctx, "SELECT * FROM pet")
  if err != nil {
    log.Println(err)
    return
  }
  defer rows.Close()

  bks := make([]*Pet, 0)
  for rows.Next() {
    bk := new(Pet)
    if err := rows.Scan(&bk.name, &bk.owner, &bk.species, &bk.sex, &bk.birth); err != nil {
      log.Println(err)
      continue
    }
    bks = append(bks, bk)
  }

  if err = rows.Err(); err != nil {
    log.Println(err)
    return
  }

  for _, bk := range bks {
    fmt.Printf("%s, %s, %s, %s, %s\n", bk.name, bk.owner, bk.species, bk.sex, bk.birth)
  }
}


func callExample(ctx context.Context) {
  db, err := sql.Open("mysql", "root:cavisson@tcp(localhost:3306)/testdb")
  if err != nil {
    log.Println("Error opening DB")
    return
  }
 

  // Set the maximum number of open connections
  db.SetMaxOpenConns(10)

  // Set the maximum number of idle connections
  db.SetMaxIdleConns(5)

  // Ping the database once to verify the connection
  if err = db.Ping(); err != nil {
    log.Println("Error connecting to DB")
    return
  }

  for i := 0; i < 10; i++ {
    pingAndQueryDB(db, ctx)
  }
}

func parseTemplates(filenames ...string) (*template.Template, error) {
  return template.ParseFiles(filenames...)
}


func main() {

  err := godotenv.Load()
  if err != nil {
      log.Println("Error loading .env file")
  }

  // Echo instance
  e := echo.New()
  mgr := statsview.New()

  // Start() runs a HTTP server at `localhost:18066` by default.
  go mgr.Start()

  // Serve static files
  staticDir := http.Dir("static")
  staticFileServer := http.FileServer(staticDir)
  e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", staticFileServer)))

  // Middleware
  /* e.Use(middleware.Logger())
  e.Use(middleware.Recover())*/
  e.Renderer = echotemplate.Default()
  templateFiles := []string{
  "view/home.html",
  "view/about.html",
  "view/page.html",
  "view/cassandra.html",
  "view/mongodb.html",
  "view/json.html",
  "view/postpg.html",
  "view/mysql.html",
  "view/cats.html",
  "view/httpcall.html",
  "view/hello.html",
  "view/database.html",
  "view/generator.html",
  "view/bost.html",
  }
  
  templates := make(map[string]*template.Template)
  for _, filename := range templateFiles {
    tmpl, err := parseTemplates(filename, "view/base.html")
    if err != nil {
      // Handle the error appropriately (e.g., log, return an error)
      log.Println("error in templates")
    }
    templates[filepath.Base(filename)] = tmpl
  }

  e.Renderer = &TemplateRegistry{
    templates: templates,
  }
  //grilla framework
  grouter := mux.NewRouter()
  grouter.HandleFunc("/mux/", handler.HttpHandler).Methods("GET", "PUT").Name("mux")

  gorillaRouterNames := map[string]string{
    "mux": "/mux/",
  }
  for name, url := range gorillaRouterNames {
    route := grouter.GetRoute(name)
    methods, _ := route.GetMethods()

    e.Match(methods, url, echo.WrapHandler(route.GetHandler()))

  }



  e.GET("/hello", Hello)
  e.GET("/QRcode", homeHandler)
  e.POST("/generator", viewCodeHandler)
  // Route => handler
  e.GET("/", handler.HomeHandler)
  e.GET("/about", handler.AboutHandler)

  // Route => api
  e.GET("/api/get-full-name", api.GetFullName)
  e.POST("/api/post-full-name", api.PostFullName)
  e.GET("/database", Database)
  e.GET("/database/mysql", mainAdmin)
  e.GET("/database/mongodb", mainmongo)
  e.GET("/database/cassandra", maincassandra)

  e.GET("/database/postpg", Postpg)
  e.GET("/hello/cats", getabhi)
  e.GET("/httpcall",mainAdmin2)
  e.POST("/save",handler.SaveDB)
  e.POST("/submit", handler.Complete)
  e.GET("/hello/json/abhi/karan/abhay", func(c echo.Context) error {
    return c.Render(http.StatusOK, "json.html", echo.Map{"title": "Page file title!!"})
  })

  e.GET("/page", func(c echo.Context) error {
    return c.Render(http.StatusOK, "page.html", echo.Map{"title": "Page file title!!"})
  })
 
  e.GET("/bost", func(c echo.Context) error {
        return c.Render(http.StatusOK, "bost.html", echo.Map{"title": "Page file title!!"})
  })

  signalChannel := make(chan os.Signal, 1)
  signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
  go func() {
    <-signalChannel
    fmt.Println("Received termination signal. Closing the database connection...")
    os.Exit(1)
  }()

  /*file, err := os.OpenFile("/home/cavisson/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
      log.Fatal(err)
    }


    log.SetOutput(file)*/

    defer func() {
      if r := recover(); r != nil {
        log.Fatal("Panic occurred:", r)
      }
      // fmt.Println("Closing the database connection...")
      //   db.Close()
    }()
    

  e.Logger.Fatal(e.Start(":1323"))
}

func homeHandler(c echo.Context) error {
  

  return c.Render(http.StatusOK, "generator.html", echo.Map{"title": "Page file title!!"})
}


func viewCodeHandler(c echo.Context) error {
  dataString := c.FormValue("dataString")

  qrCode, _ := qr.Encode(dataString, qr.L, qr.Auto)
  qrCode, _ = barcode.Scale(qrCode, 200, 200)

  buffer := new(bytes.Buffer)
  png.Encode(buffer, qrCode)

  return c.Blob(http.StatusOK, "image/png", buffer.Bytes())
}

func mainAdmin2(c echo.Context) error {
  return c.Render(http.StatusOK, "httpcall.html", echo.Map{"title": "Page file title!!"})
}

func Database(c echo.Context) error {

  return c.Render(http.StatusOK, "database.html", echo.Map{"title": "Page file title!!"})
}


func Hello(c echo.Context) error { 
  return c.Render(http.StatusOK, "hello.html", echo.Map{"title": "Page file title!!"})
}

func getabhi(c echo.Context) error {

  catName := c.QueryParam("name")
  catType := c.QueryParam("type")
  
  log.Printf("\nname=%s,type=%s", catName, catType)
  return c.Render(http.StatusOK, "cats.html", echo.Map{"title": "Page file title!!"})
}
func mainAdmin(c echo.Context) error {
  
  log.Println("Hello mainAdmin!")
  req := c.Request()
  ctx := req.Context()
  callExample(ctx)
  return c.Render(http.StatusOK, "mysql.html", echo.Map{"title": "Page file title!!"})

}


func mainmongo(c echo.Context) error {
    req := c.Request()
    ctx := req.Context()

    // Set client options
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

    // Connect to MongoDB
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Println(err)
        return err
    }
    defer client.Disconnect(ctx)

    // Ping the MongoDB server
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Println(err)
        return err
    }

    // Access the collection
    collection := client.Database("test").Collection("trainers")

    ash := Trainer{"Ash", 10, "Pallet Town"}

    _, err = collection.InsertOne(ctx, ash)
    if err != nil {
        log.Println(err)
        return err
    }

    filter := bson.D{{"name", "Ash"}}

    var result Trainer

    err = collection.FindOne(ctx, filter).Decode(&result)
    if err != nil {
        log.Println(err)
        return err
    }

    // testinstrument()

    return c.Render(http.StatusOK, "mongodb.html", echo.Map{"title": "Page file title!!"})
}

/*func testinstrument() {
    log.Println("inside testinstrument")
    time.Sleep(1 * time.Second)
}
*/

func maincassandra(c echo.Context) error {
    req := c.Request()
    ctx := req.Context()

    callGOCql(ctx)

    return c.Render(http.StatusOK, "cassandra.html", echo.Map{"title": "Page file title!!"})
}

type Emp struct {
    id        string
    firstName string
    lastName  string
    age       int
}

func createEmp(session *gocql.Session, emp Emp) error {
    query := "INSERT INTO emps(empid, first_name, last_name, age) VALUES(?, ?, ?, ?)"
    err := session.Query(query, emp.id, emp.firstName, emp.lastName, emp.age).Exec()
    if err != nil {
        return err
    }
    return nil
}

func callGOCql(ctx context.Context) {
    cluster := gocql.NewCluster("127.0.0.1")
    cluster.Keyspace = "code2succeed"
    session, err := cluster.CreateSession()
    if err != nil {
        log.Println("Cassandra not connected")
        return
    }
    defer session.Close()

    emp1 := Emp{"E3451", "Anupadfm", "Rsdaj", 21}
    emp2 := Emp{"E45672", "Rahdul", "Anadfnd", 30}

    if err := createEmp(session, emp1); err != nil {
        log.Println(err)
    }

    if err := createEmp(session, emp2); err != nil {
        log.Println(err)
    }
}

func checkErr(err error) {
  if err != nil {
    log.Println("Not find")
  }
}

func Postpg(c echo.Context) error {
    dbUser, dbPassword := os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS")
    if dbUser == "" && dbPassword == "" {
        dbUser = "postgres"
        dbPassword = "12345678"
    }

    options := &pg.Options{
        User:     dbUser,
        Password: dbPassword,
        Database: "postgres",
    }
    db := pg.Connect(options)
    defer db.Close()

    var students []Student
    err := db.Model(&students).Select()
    if err != nil {
        log.Println(err)
        return err
    }

    var student Student
    err = db.Model(&student).Where("Rollno = ?", 1).Select()
    if err != nil {
        log.Println(err)
        return err
    }

    /*fmt.Println(students)
    fmt.Println(student)*/

    return c.Render(http.StatusOK, "postpg.html", echo.Map{"title": "Page file title!!"})
}
