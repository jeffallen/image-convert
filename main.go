package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"strings"

	// Load all the file formats
	_ "code.google.com/p/go.image/bmp"
	_ "code.google.com/p/go.image/tiff"
	_ "github.com/ftrvxmtrx/groke/image/pcx"
	_ "github.com/jeffallen/g/image/xbm"
	_ "github.com/knieriem/g/image/pnm"
	_ "github.com/mewrnd/blizzconv/images/cel"
	_ "github.com/mewrnd/blizzconv/images/cl2"
	_ "image/gif"
	_ "image/jpeg"
)

func convert(w http.ResponseWriter, r *http.Request) {
	f, fh, err := r.FormFile("img")
	if err != nil {
		http.Error(w, "Missing image data.", 500)
		return
	}

	img, kind, err := image.Decode(f)
	if err != nil || img == nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fn := fh.Filename
	dot := strings.LastIndex(fn, ".")
	if dot >= 0 {
		fn = fn[0:dot]
	}
	fn = fmt.Sprintf("%v.png", fn)

	log.Print("Image filename: ", fh.Filename,
		" type: ", fh.Header["Content-Type"],
		" decoded-as: ", kind,
		" new-fn: ", fn)

	w.Header().Add("Content-Disposition",
		fmt.Sprintf("attachment; filename=\"%v\"", fn))
	err = png.Encode(w, img)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	return
}

func main() {
	log.Println("Starting with args:", os.Args)
	http.HandleFunc("/convert", convert)
	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
