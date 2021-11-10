package main

import (
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"

	"bigchief64/territoryTool/csv"

	"time"

	"github.com/harry1453/go-common-file-dialog/cfd"
	"github.com/zserge/lorca"
)

//go:embed views
var fs embed.FS
var debug = false

type interact struct {
	sync.Mutex
	ui lorca.UI
}

func (c *interact) Close() {
	c.Lock()
	defer c.Unlock()
	c.ui.Close()
	os.Exit(3)
}

func (c *interact) OpenFile() string {
	c.Lock()
	defer c.Unlock()

	f := openDialog()
	data := csv.GetData(f)

	//loop through each line and clean
	for k, v := range data {
		if strings.Contains(v[2], "Ponchatoula") {
			i := strings.LastIndex(v[2], "Ponchatoula")
			data[k][2] = v[2][:i]
			data[k] = append(v, "Ponchatoula, LA 70454")
		} else if strings.Contains(v[2], "Springfield") {
			i := strings.LastIndex(v[2], "Springfield")
			data[k][2] = v[2][:i]
			data[k] = append(v, "Springfield, LA 70462")
		} else if strings.Contains(v[2], "Hammond") {
			i := strings.LastIndex(v[2], "Hammond")
			data[k][2] = v[2][:i]
			data[k] = append(v, "Hammond, LA 70403")
		}
		tempStr := convertStrHTML(data[k])
		tempStr = convertToLink(tempStr, data[k][2])
		data[k] = append(data[k], tempStr)
	}

	return formatIntoHTML(data)
}

func convertToLink(str, name string) string{
	tempStr := "<a href='" + str + "' target='_blank'>" + name + "</a>"

	return tempStr
}

func convertStrHTML(lines []string) string {
	str := "https://fastpeoplesearch.com/address/"
	var str2, str3 string

	if len(lines[2]) > 0 {
		str2 = lines[2][:len(lines[2])-1]
		str2 = strings.ReplaceAll(str2, " ", "-")
	}

	if len(lines) > 3 {
		str3 = strings.ReplaceAll(lines[3][:len(lines[3])-6], " ", "-")
		str3 = strings.Replace(str3, ",", "", -1)
		str = str + str2 + "_" + str3
		return str
	}

	return str2
}

func formatIntoHTML(data map[string][]string) string {
	var dataString string
	dataString = "<table><thead><tr><th>Assigned</th><th>Name</th><th>Sex</th>" + 
		"<th>Phone</th><th>Address</th><th>City</th><th>Age</th><th>Link</th></tr></thead>"

	for k, v := range data {
		dataString = dataString + "<tr><td></td><td>" + k + "</td>"
		for i, r := range v {
			dataString = dataString + "<td>" + r + "</td>"
			if i == 3 {
				dataString = dataString + "<td>     </td>"
			}
		}
		dataString = dataString + "</tr>"
	}
	dataString = dataString + "</table>"

	return dataString
}

func main() {
	args := []string{""}
	if debug {
		args = append(args, "debug")
	} else {
		args = append(args, "")
	}
	args = append(args, "")
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	ui, err := lorca.New("", "", 1024, 768, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	// A simple way to know when UI is ready (uses body.onload event in JS)
	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

	//enable exit function
	c := &interact{}
	c.ui = ui
	ui.Bind("close", c.Close)
	ui.Bind("openFile", c.OpenFile)

	// Load HTML.
	// You may also use `data:text/html,<base64>` approach to load initial HTML,
	// e.g: ui.Load("data:text/html," + url.PathEscape(html))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	go http.Serve(ln, http.FileServer(http.FS(fs)))
	//load starting page
	if debug {
		ui.Load(fmt.Sprintf("http://%s/views/index.html", ln.Addr()))
	} else {
		ui.Load(fmt.Sprintf("http://%s/views/index.html", ln.Addr()))
	}

	// You may use console.log to debug your JS code, it will be printed via
	// log.Println(). Also exceptions are printed in a similar manner.
	ui.Eval(`
		console.log("Hello console");
		console.log('Multiple values:', [1, false, {"x":5}]);
	`)

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal, 4)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:

	case <-ui.Done():
	}

	log.Println("exiting...")
}

func openDialog() string {
	openDialog, err := cfd.NewOpenFileDialog(cfd.DialogConfig{
		Title: "Open A File",
		Role:  "OpenFileExample",
		FileFilters: []cfd.FileFilter{
			{
				DisplayName: "CSV Files (*.csv)",
				Pattern:     "*.csv",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
		SelectedFileFilterIndex: 0,
		FileName:                "file.csv",
		DefaultExtension:        "csv",
	})
	if err != nil {
		log.Fatal(err)
	}
	//this was added as a quick fix to a race condition
	go func() {
		time.Sleep(2 * time.Second)
		if err := openDialog.SetFileName("hello world"); err != nil {
			log.Fatal(err)
		}
	}()
	if err := openDialog.Show(); err != nil {
		log.Fatal(err)
	}
	result, err := openDialog.GetResult()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Chosen file: %s\n", result)
	return result
}
