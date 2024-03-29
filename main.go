package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const COMBINED_FILE_NAME = "combined.csv"
const WORDS_FOLDER_NAME = "./word-list"

var normalizer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

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
	totalUniqueWords := 0

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
		totalUniqueWords = totalUniqueWords + uniqueWords

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
	fmt.Println("File written successfully. Please see the file named ", COMBINED_FILE_NAME, "in the ", WORDS_FOLDER_NAME, "folder to see the combined list. Total parsed word ", totalUniqueWords)

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

		// parse the data
		word := strings.TrimSpace(strings.ToLower(data[0]))
		// do some processing to remove non printable unicode character like ZERO WIDTH NO-BREAK SPACE (\uFEFF)
		processedWord := strings.Map(func(r rune) rune {
			if unicode.IsPrint(r) {
				return r
			}
			return -1
		}, word)

		// normalize word from accent to english word
		processedWordWithBrackets := processWord(processedWord)

		if len(processedWordWithBrackets) > 0 {
			for _, word := range processedWordWithBrackets {
				fileData = append(fileData, []string{word})
			}

		}

	}

	return fileData, totalLine
}

// func normalize(str string) (string, error) {
// 	s, _, err := transform.String(normalizer, str)
// 	if err != nil {
// 		return "", err
// 	}
// 	return strings.ToLower(s), err
// }

func processWord(str string) []string {

	str = strings.TrimSpace(strings.Join(strings.Fields(str), " "))

	s, _, err := transform.String(normalizer, str)

	if err != nil {
		return []string{}
	}

	str = strings.ToLower(s)
	// replace underscores and slashes
	processedWord := strings.Replace(str, "'", "", -1)
	processedWord = strings.Replace(processedWord, "\\", "", -1)
	processedWord = strings.Replace(processedWord, "_", "-", -1)

	matchWords := match(processedWord)

	return matchWords
}

func match(s string) []string {
	str := []string{}
	i := strings.Index(s, "(")
	if i >= 0 {
		j := strings.Index(s, ")")
		if j >= 0 {
			//return s[i+1 : j]
			str = append(str, s[i+1:j])
			str = append(str, s[:i])
		}
	} else {
		str = append(str, s)
	}

	return str
}
