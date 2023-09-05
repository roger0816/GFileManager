package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	http.HandleFunc("/", uploadFile)
	http.HandleFunc("/download/", downloadFile)
	http.ListenAndServe(":8080", nil)
}

func getGoFilePath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "./"
	}
	fmt.Printf("path : %s\n", dir)
	return dir
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(10 << 20) // limit to 10MB

		file, fileHeader, err := r.FormFile("fileToUpload")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}
		defer file.Close()

		savePath := filepath.Join(getGoFilePath(), fileHeader.Filename)
		tempFile, err := os.Create(savePath)
		if err != nil {
			fmt.Println("Error Creating File")
			fmt.Println(err)
			return
		}
		defer tempFile.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			fmt.Println("Error Reading File")
			fmt.Println(err)
			return
		}
		tempFile.Write(fileBytes)
		fmt.Fprintf(w, "Successfully Uploaded File to %s\n", savePath)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`
        <form enctype="multipart/form-data" action="/" method="post">
            <input type="file" name="fileToUpload">
            <input type="submit" value="Upload">
        </form>
		<br>
		<a href="/download/">Click here to download</a>
    `))
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Filename not specified", http.StatusBadRequest)
		return
	}

	filepath := filepath.Join(getGoFilePath(), filename)

	// Set the Content-Type based on file extension
	//ext := filepath.Ext(filename)
	ext := getExt(filename)
	switch ext {
	case ".exe":
		w.Header().Set("Content-Type", "application/octet-stream")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".zip":
		w.Header().Set("Content-Type", "application/zip")
	// Add more cases if needed
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	http.ServeFile(w, r, filepath)
}

func getExt(fileName string) string {
	for i := len(fileName) - 1; i >= 0 && !os.IsPathSeparator(fileName[i]); i-- {
		if fileName[i] == '.' {
			return fileName[i:]
		}
	}
	return ""
}
