package main

import (
	"bufio"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"iCASComaasJoelPintoMata/utils"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"gopkg.in/russross/blackfriday.v2"
)

const csvFileType string = "csv"
const prnFileType string = "prn"

var regexHeader *regexp.Regexp
var regexBody *regexp.Regexp

// initializes this process
func init() {
	// lets compile the regular expressions just once
	patternHeader := `(?P<name>(.*?))\,(?P<address>(.*?))\,(?P<postcode>(.*?))\,(?P<phonenumber>(.*?))\,(?P<creditLimit>(.*?))\,(?P<birthday>(.*))`
	thisRegexHeader, err := regexp.Compile(patternHeader)
	if err != nil {
		log.Fatal(err)
	}
	regexHeader = thisRegexHeader

	patternBody := `(?P<name>\"(.*?)\")\,(?P<address>(.*?))\,(?P<postcode>(.*?))\,(?P<phonenumber>(.*?))\,(?P<creditLimit>(.*?))\,(?P<birthday>(.*))`
	thisRegexBody, err := regexp.Compile(patternBody)
	if err != nil {
		log.Fatal(err)
	}
	regexBody = thisRegexBody
}

func main() {
	// starting gin server
	router := gin.Default()

	router.GET("/convert/"+csvFileType, convertCSVtoHTML)
	router.GET("/convert/"+prnFileType, convertPRNtoHTML)

	// Register pprof handlers
	pprof.Register(router, &pprof.Options{
		// default is "debug/pprof"
		RoutePrefix: "debug/pprof",
	})

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}

// Auxiliary function orchetrating the csv to html convertion
func convertCSVtoHTML(c *gin.Context) {
	getHTML(c, csvFileType, "Workbook2.csv")
}

// Auxiliary function orchetrating the prn to html convertion
func convertPRNtoHTML(c *gin.Context) {
	getHTML(c, prnFileType, "Workbook2.prn")
}

// Reads a file [CSV|PRN] and convert its contents into HTML code
func getHTML(c *gin.Context, fileType string, fileName string) {
	var html string
	var headerColumnIndexArray []int

	isHeader := true

	// Open file and create scanner on top of it
	file, err := os.Open(fileName)
	if err != nil {
		setError(c, err)
		return
	}
	defer file.Close()

	// we detected csv characters encoded with ISO8859_1
	// as in http://stackoverflow.com/questions/29686673/read-unicode-characters-with-bufio-scanner-in-go
	decoded := transform.NewReader(file, charmap.ISO8859_1.NewDecoder())
	scanner := bufio.NewScanner(decoded)

	for scanner.Scan() {

		string := scanner.Text()

		// get the header in order to dynamically determine the header indexes
		if isHeader {
			headerColumnIndexArray = getHeaderColumnIndex(string)
			string, err = getTableHeader(fileType, headerColumnIndexArray, string)
			isHeader = false
		} else {
			string, err = getTableLine(fileType, headerColumnIndexArray, string)
		}

		// check for errors found while processing the string
		if err != nil {
			setError(c, err)
			return
		}
		html += string
	}

	html = getTableWrap(html)

	markdown := blackfriday.Run([]byte(html))
	c.Data(http.StatusOK, "text/html; charset=utf-8", markdown)
}

/**
Processes an application error
1) Logs the error
2) Sets the error into the gin context
*/
func setError(c *gin.Context, err error) {
	markdown := blackfriday.Run([]byte(err.Error()))
	c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", markdown)
	log.Panic(err)
}

// Performs a final table wrap around the given html code
func getTableWrap(html string) string {
	return "<table border='1'>" + html + "</table>"
}

// Performs a final table row wrap around the given html code
func getRowWrap(html string) string {
	return "<tr>" + html + "</tr>"
}

// Performs a final table cell wrap around the given html code
func getCellWrap(html string, align string) string {
	if strings.EqualFold(align, "right") {
		html = "<td align='right'><pre>" + html + "</pre></td>"
	} else {
		html = "<td><pre>" + html + "</pre></td>"
	}
	return html
}

// Calculates the index position(s) of the header column(s) on a given PRN type string
func getHeaderColumnIndex(header string) []int {
	return []int{
		strings.Index(header, "Birthday"),
		strings.Index(header, "Credit Limit"),
		strings.Index(header, "Phone"),
		strings.Index(header, "Postcode"),
		strings.Index(header, "Address"),
		strings.Index(header, "Name")}
}

// Builds an html table header line
func getTableHeader(fileType string, headerColumnIndexArray []int, line string) (string, error) {
	var value string
	var err error

	if strings.EqualFold(fileType, csvFileType) {
		value, err = getCSVToTableRow(regexHeader, line)
	}

	if strings.EqualFold(fileType, prnFileType) {
		value = getPRNToTableRow(headerColumnIndexArray, line)
	}
	return value, err
}

// Builds an html table header line
func getTableLine(fileType string, headerColumnIndexArray []int, line string) (string, error) {
	var value string
	var err error

	if strings.EqualFold(fileType, csvFileType) {
		value, err = getCSVToTableRow(regexBody, line)
	}

	if strings.EqualFold(fileType, prnFileType) {
		value = getPRNToTableRow(headerColumnIndexArray, line)
	}

	return value, err
}

// Parses a CSV type row into a html string
// check regex at https://regex-golang.appspot.com/assets/html/index.html
func getCSVToTableRow(regexp *regexp.Regexp, line string) (string, error) {
	var html string
	var err error

	matches := regexp.FindStringSubmatch(line)
	if matches == nil {
		html = line
		err = errors.New("Error while parsing the line (check the require format in the documentation)")
	} else {
		// extract the name
		html += getCellWrap(matches[2], "left")
		// extract the address
		html += getCellWrap(matches[3], "left")
		// extract the postcode
		html += getCellWrap(matches[5], "left")
		// extract the phone number
		html += getCellWrap(matches[7], "left")
		// extract the credit limit
		html += getCellWrap(matches[9], "right")
		// extract the birthday
		html += getCellWrap(matches[11], "left")

		html = getRowWrap(html)
	}
	return html, err
}

// Parses a PRN type row into a html string
func getPRNToTableRow(headerColumnIndexArray []int, line string) string {
	var beginCell = "<td><pre>"
	var endCell = "</pre></td>"

	for _, element := range headerColumnIndexArray {
		if element > 0 {
			line = utils.GetChars(line, 0, element-1) + endCell + beginCell + utils.GetChars(line, element, len([]rune(line)))
		}
	}
	return getRowWrap(beginCell + line + endCell)
}
