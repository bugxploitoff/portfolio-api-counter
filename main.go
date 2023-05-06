package main

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "github.com/rs/cors"
  "os"
)

var views = 0

type ViewsData struct {
    Views int `json:"views"`
}

type ContactData struct {
    Email   string `json:"email"`
    Message string `json:"message"`
}

func main() {
    // Load the previous value of views from the file
    data, err := ioutil.ReadFile("views.json")
    if err == nil {
        var viewsData ViewsData
        err := json.Unmarshal(data, &viewsData)
        if err == nil {
            views = viewsData.Views
        }
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/api/portfolio/views", viewsHandler)
    mux.HandleFunc("/api/portfolio/contact", contactHandler)
    mux.Handle("/", http.FileServer(http.Dir("public")))

    handler := cors.Default().Handler(mux)
    http.ListenAndServe(":80", handler)
}

func viewsHandler(w http.ResponseWriter, r *http.Request) {
    // Increment the views count
    views++

    // Store the new value in the file
    viewsData := ViewsData{Views: views}
    jsonData, err := json.Marshal(viewsData)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    err = ioutil.WriteFile("views.json", jsonData, 0644)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    data := map[string]int{"views": views}
    json.NewEncoder(w).Encode(data)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
    // Read the request body
    requestBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Parse the JSON request body into ContactData struct
    var contactData ContactData
    err = json.Unmarshal(requestBody, &contactData)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Append the new contact data to the file
    var contacts []ContactData

    file, err := ioutil.ReadFile("contacts.json")
    if err != nil {
        if !os.IsNotExist(err) {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        file = []byte("[]")
    }

    err = json.Unmarshal(file, &contacts)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    contacts = append(contacts, contactData)

    jsonData, err := json.Marshal(contacts)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    err = ioutil.WriteFile("contacts.json", jsonData, 0644)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Respond with success message
    response := map[string]string{
        "status":  "success",
        "message": "Contact data stored successfully",
    }
    json.NewEncoder(w).Encode(response)
}

