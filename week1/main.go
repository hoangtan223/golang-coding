package main

import (
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	configureCommand := flag.NewFlagSet("configure", flag.ExitOnError)
	address := configureCommand.String("a", "", "Shorten Address")
	targetUrl := configureCommand.String("u", "", "Target URL")

	runCommand := flag.NewFlagSet("run", flag.ExitOnError)
	runPort := runCommand.Int("p", -1, "Running on Port")
	deleteTarget := flag.String("d", "", "Delete Target Address")
	listTargets := flag.Bool("l", false, "List shorten Url list")
	printUsage := flag.Bool("h", false, "Print Usage")

	if *printUsage {
		flag.Usage()
		return
	}

	switch os.Args[1] {
	case "configure":
		configureCommand.Parse(os.Args[2:])
	case "run":
		runCommand.Parse(os.Args[2:])
	default:
		flag.Parse()
	}
	mapList := GetList()

	if *listTargets {
		for addressVal, targetVal := range mapList {
			fmt.Printf("%s: %s\n", addressVal, targetVal)
		}
		return
	}

	if *address != "" && *targetUrl != "" {
		mapList[*address] = *targetUrl
		WriteYAML(mapList)
		return
	}

	if *deleteTarget != "" {
		DeleteAdress(*deleteTarget)
		return
	}

	if *runPort != -1 {
		http.HandleFunc("/", redirect)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func GetList() map[string]string {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var mapList map[string]string
	err = yaml.Unmarshal(yamlFile, &mapList)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return mapList
}

func WriteYAML(yamlStr map[string]string) {
	data, _ := yaml.Marshal(yamlStr)
	err := ioutil.WriteFile("config.yaml", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteAdress(address string) {
	ymlContent := GetList()
	_, ok := ymlContent[address]
	if ok {
		delete(ymlContent, address)
		WriteYAML(ymlContent)
	} else {
		fmt.Println("Address not exist")
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	path := html.EscapeString(r.URL.Path)
	if r.URL.Path == "/" {
		fmt.Fprintf(w, "Hi there, welcome to urlshortener")
	} else {
		list := GetList()
		http.Redirect(w, r, list[path[1:]], 301)
	}
}
