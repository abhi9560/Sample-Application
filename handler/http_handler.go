package handler

import (
    "fmt"
    /*"github.com/gorilla/mux"*/
    "net/http"
)

const AddForm = `<!DOCTYPE html>
<html lang="en">
    <style>
        .hello{ color: red;}
        hr{ border: 1px #ccc dashed;}
    </style>


<head>
    <meta charset="UTF-8">
    <meta http-equiv="x-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width", inital-scale="1.0">
    <title>Home page</title>
    <link rel="stylesheet" type="text/css" href="/static/css/main.css">
</head>
<body>  
    <header>
    <img src="/static/image/goimg.png" alt ="Gopher">
    <b style='color:red;'>Cavisson System </b>
    </header>
    <nav>
    <ul>           
        <li style='color::wred;'><a href="/"; style='color:lightyellow;'> Home </a></li>
        <li style='color:red;'><a href="/page"; style='color:lightyellow;'>Content </a></li>
        <li style='color:red;'><a href="/about"; style='color:lightyellow;'>About</a></li>
    </ul>
    </nav> 
</body>
</html>`

func HttpHandler(w http.ResponseWriter, r *http.Request) {

    fmt.Fprintf(w, AddForm)
    fmt.Fprintf(w, "<h1 style='color:red;'>Hello Gorilla</h1>")

    w.Header().Set("Content-Type", "application/json")
    // Please note the the second parameter "home.html" is the template name and should
    // be equal to one of the keys in the TemplateRegistry array defined in main.go

}

/*func HttpHandlerP(w http.ResponseWriter, r *http.Request) {

    fmt.Fprintf(w, AddForm21)
    w.Header().Set("Content-Type", "application/json")
    // Please note the the second parameter "home.html" is the template name and should
    // be equal to one of the keys in the TemplateRegistry array defined in main.go

}*/
