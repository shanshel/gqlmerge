package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type input struct {
	src  string
	dest string
}

func main() {

	flag.Usage = func() {
		fmt.Println("Usage: gqlmerge $srcFolder $destFile")
		flag.PrintDefaults()
	}

	//Get CLI Inputs
	inputData, err := getCliInput()
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(inputData.src)
	if err != nil {
		log.Fatal(err)
	}

	//Define Regexes
	queryRegex, err := regexp.Compile("type[\\s\\n]+Query[\\s\\n]+\\{[\\s\\S\n]+?\\}")
	if err != nil {
		log.Fatal(err)
	}

	mutationRegex, err := regexp.Compile("type[\\s\\n]+Mutation[\\s\\n]+\\{[\\s\\S\n]+?\\}")
	if err != nil {
		log.Fatal(err)
	}

	subscriptionRegex, err := regexp.Compile("type[\\s\\n]+Subscription[\\s\\n]+\\{[\\s\\S\n]+?\\}")
	if err != nil {
		log.Fatal(err)
	}

	inBetween, err := regexp.Compile("\\{[\\s\\S\\n]+\\}")
	if err != nil {
		log.Fatal(err)
	}

	container := ""
	//Define container for speical graphql types
	queryContainer := "\n\rtype Query {\n\r"
	mutationContainer := "\n\rtype Mutation {\n\r"
	subscriptionContainer := "\n\rtype Subscription {\n\r"
	for _, file := range files {

		ext := filepath.Ext(file.Name())
		if ext != ".graphqls" {
			continue
		}

		fileContentBytes, err := os.ReadFile(path.Clean(inputData.src + "/" + file.Name()))
		if err != nil {
			log.Fatal(err)
		}

		fileContentString := string(fileContentBytes)

		//Find Graphql Types
		fQuery := queryRegex.FindString(fileContentString)
		fMutation := mutationRegex.FindString(fileContentString)
		fSubscription := subscriptionRegex.FindString(fileContentString)

		if fQuery != "" {
			fileContentString = strings.Replace(fileContentString, fQuery, "", -1)
			fQuery = inBetween.FindString(fQuery)

			queryContainer += fQuery[1 : len(fQuery)-1]

		}

		if fMutation != "" {
			fileContentString = strings.Replace(fileContentString, fMutation, "", -1)
			fMutation = inBetween.FindString(fMutation)
			mutationContainer += fMutation[1 : len(fMutation)-1]
		}

		if fSubscription != "" {
			fileContentString = strings.Replace(fileContentString, fSubscription, "", -1)
			fSubscription = inBetween.FindString(fSubscription)
			subscriptionContainer += fSubscription[1 : len(fSubscription)-1]

		}
		container += "\n#_________________ " + file.Name() + " _________________\n"
		container += strings.TrimSpace(fileContentString)

	}
	queryContainer += " \n\r}"
	mutationContainer += " \n\r}"
	subscriptionContainer += " \n\r}"
	container += "\n#_________________ GraphQL Types _________________\n"

	container += queryContainer + mutationContainer + subscriptionContainer

	os.WriteFile(inputData.dest, []byte(strings.TrimSpace(container)), 0644)
}

func getCliInput() (*input, error) {
	if len(os.Args) < 2 {
		return nil, errors.New("src & dist argument are required, example: gogracom $src $dest")
	}
	flag.Parse()
	return &input{
		src:  flag.Arg(0),
		dest: flag.Arg(1),
	}, nil
}
