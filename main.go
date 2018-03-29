package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles("template/upload.html"))

//show the named template
func show(w http.ResponseWriter, template string, data interface{}) {
	templates.ExecuteTemplate(w, template+".html", data)
}

//This is where the action happens.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//GET shows the Upload Form
	case "GET":
		show(w, "upload", nil)

	//POST stores the uploaded files into uploads folder.
	case "POST":
		//Multipart form parsing.
		err := r.ParseMultipartForm(100000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//get a ref to the parsed multipart form
		m := r.MultipartForm

		//get the *fileheaders
		files := m.File["myfiles"]
		for i := range files {
			//for each fileheader, get a handle to the actual file
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//Destination file creation.
			dst, err := os.Create("./uploads/" + files[i].Filename)
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//Copy the uploaded file content to the dest file.
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		//Show success message.
		show(w, "upload", "Upload successful.")

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/upload", uploadHandler)

	//static file handler.
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))

	//Listen on port 2023
	http.ListenAndServe(":2023", nil)
}
