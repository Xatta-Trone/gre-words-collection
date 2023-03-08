package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

const COMBINED_FILE_NAME = "combined.csv"
const WORDS_FOLDER_NAME = "./word-list"

func main() {

	fmt.Println("Welcome")

	var filenames []string

	// get file list from directory
	files, err := os.ReadDir(WORDS_FOLDER_NAME)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		// only add the files other than the combined file
		if file.Name() != COMBINED_FILE_NAME {
			filenames = append(filenames, WORDS_FOLDER_NAME+"/"+file.Name())
		}

	}

	// read the files
	// then filter the words by looking up into to word key function
	wordKeys := make(map[string]bool)
	combinedList := [][]string{}

	for _, fileName := range filenames {

		fmt.Println("Reading " + fileName)

		fileData, total := readCSV(fileName)

		uniqueWords := 0

		for _, item := range fileData {
			if _, value := wordKeys[item[0]]; !value {
				wordKeys[item[0]] = true
				combinedList = append(combinedList, item)
				uniqueWords++
			}
		}

		fmt.Println(fileName+"has "+strconv.Itoa(total)+" words; among them ", strconv.Itoa(uniqueWords)+" unique words found")
		fmt.Println("==========================================================================================================")

	}

	// sort the words alphabetically
	sort.Slice(combinedList, func(i, j int) bool {
		return combinedList[i][0] < combinedList[j][0]
	})

	// now write the combinedList into a csv
	combinedFile, err := os.Create(WORDS_FOLDER_NAME + "/" + COMBINED_FILE_NAME)

	if err != nil {
		log.Fatal("failed to open file", err)
	}

	defer combinedFile.Close()

	w := csv.NewWriter(combinedFile)
	err = w.WriteAll(combinedList)

	if err != nil {
		log.Fatal("Could not write to combined file", err)
	}

	// fmt.Println(combinedList)
	fmt.Println("File written successfully. Please see the file named ", COMBINED_FILE_NAME, "in the ", WORDS_FOLDER_NAME, "folder to see the combined list.")

}

func readCSV(filePath string) ([][]string, int) {
	var totalLine int
	var fileData [][]string
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Could not open the file", filePath, err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)

	for {
		data, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// do something with read line
		// count the number of rows
		totalLine++
		// add data to temp array
		fileData = append(fileData, data)
	}

	return fileData, totalLine
}
