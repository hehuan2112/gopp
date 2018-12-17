package main

import (
	"os"
    "log"
	"fmt"
	"strings"
    "html/template"
	"net/http"
	"io/ioutil"
	"path/filepath"
)

const tpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Go Pic Present</title>
    <style>
    html, body { width: 100%; margin: 0; padding: 0; }
    #header { position: absolute; top: 0; left: 0; width: 100%; height: 24px; 
              margin: 0 0 10px 0; border-bottom: 1px solid #EAEAEA; }
    #myh1 { float: left; width: 100px; font-size: 16px; line-height: 24px; }
    #helper { float: left; line-height: 24px; }
    #gallery { width: 100%; margin: 34px 0 0 0; }
    .pic { float: left; margin: 0 10px 10px 0; padding: 2px; border: 1px solid #EAEAEA; }
    .pic .img-1-2 img { max-width: 100px; max-height: 200px; }
    .pic .img-1-3 img { max-width: 100px; max-height: 300px; }
    .pic .img-1-4 img { max-width: 100px; max-height: 400px; }
    .pic .img-2-2 img { max-width: 200px; max-height: 200px; }
    .pic .img-2-3 img { max-width: 200px; max-height: 300px; }
    .pic .img-2-4 img { max-width: 200px; max-height: 400px; }
    .pic .img-3-2 img { max-width: 300px; max-height: 200px; }
    .pic .img-3-3 img { max-width: 300px; max-height: 300px; }
    .pic .img-3-4 img { max-width: 300px; max-height: 400px; }
    .pic .txt { color: #EAEAEA; text-align: center; }
    .pic:hover { border: 1px solid #CCCCCC; }
    .pic:hover .txt { color: black; }
    </style>
</head>
<body>
<div id="header">
    <div id="myh1"><b>@^_^@</b> GoPP</div>
    <div id="helper">
        | <input type="checkbox" id="ckb_showfn" checked="checked" onclick="toggle_filename();" /> Toggle Filename
        | <button onclick="inc_image_width();">+W</button> - <button onclick="inc_image_height();">+H</button>
    </div>
<script>
function toggle_filename() {
    console.log('* clicked toggle filename!');
    var is_checked = document.getElementById('ckb_showfn').checked;
    var fns = document.getElementsByClassName('txt');
    if (is_checked) {
        for (var i=0; i<fns.length; i++) {
            fns[i].style.display = '';
        }
    } else {
        for (var i=0; i<fns.length; i++) {
            fns[i].style.display = 'none';
        }
    }
}

var image_ws = [1, 2, 3];
var image_hs = [2, 3, 4];
var image_widx = 1;
var image_hidx = 1;

function update_image_class(class_name) {
    var imgs = document.getElementsByClassName('img');
    for (var i=0; i<imgs.length; i++) {
        var img = imgs[i];
        img.className = class_name;
    }
}

function inc_image_width() {
    image_widx = (image_widx + 1) % image_ws.length;
    image_w = image_ws[image_widx];
    image_h = image_hs[image_hidx];

    update_image_class('img img-' + image_w + '-' + image_h);
}

function inc_image_height() {
    image_hidx = (image_hidx + 1) % image_hs.length;
    image_w = image_ws[image_widx];
    image_h = image_hs[image_hidx];

    update_image_class('img img-' + image_w + '-' + image_h);
}

function dec_image_width() {
}


function dec_image_height() {
}

</script>
</div>
<div id="gallery">
{{ range .Pics }}
    <div class="pic">
        <div class="img img-2-2">
            <img src="./{{ . }}"/>
        </div>
        <div class="txt">
            {{ . }}
        </div>
    </div>
{{ end }}
</div>
</body>
</html>`


func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("* request path: ", r.URL.Path)

	if r.URL.Path == "/" {
		var pics []string
		err := filepath.Walk("./", func(path string, f os.FileInfo, err error) error {
			if f == nil { return err }
			if f.IsDir() { return nil }
			if strings.HasSuffix(f.Name(), ".png") {
				pics = append(pics, f.Name())
			}
			return nil
        })
        if err != nil {
            fmt.Printf("* filepath.Walk() returned %v\n", err)
		}
		fmt.Println("* pics: ", pics)

		data := struct {
			numPics int
			Pics []string
		}{ numPics: len(pics), Pics: pics }

		t, err := template.New("index").Parse(tpl)
		if (err != nil) {
			log.Println(err)
		}
		t.Execute(w, data)

	} else {
		file, err := os.Open("." + r.URL.Path)
		if  err != nil {
			w.Write([]byte(err.Error()))
		}

		defer file.Close()
		buff, err := ioutil.ReadAll(file)
		if  err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Write(buff)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)

    fmt.Println("* started gopp server here on 0.0.0.0:8000")
    err := http.ListenAndServe(":8000", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
