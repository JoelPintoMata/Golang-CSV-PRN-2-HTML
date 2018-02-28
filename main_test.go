package main

import (
	"regexp"
	"strings"
	"testing"

	"iCASComaasJoelPintoMata/utils"
)

/**
test suite for the CSV types to HTML converter
*/
func TestCSV(t *testing.T) {
	patternHeader := `(?P<name>(.*?))\,(?P<address>(.*?))\,(?P<postcode>(.*?))\,(?P<phonenumber>(.*?))\,(?P<creditLimit>(.*?))\,(?P<birthday>(.*))`
	thisRegexHeader, _ := regexp.Compile(patternHeader)
	regexHeader = thisRegexHeader

	patternBody := `(?P<name>\"(.*?)\")\,(?P<address>(.*?))\,(?P<postcode>(.*?))\,(?P<phonenumber>(.*?))\,(?P<creditLimit>(.*?))\,(?P<birthday>(.*))`
	thisRegexBody, _ := regexp.Compile(patternBody)
	regexBody = thisRegexBody

	// table csv header tests
	result, _ := getCSVToTableRow(regexHeader, "Name,Address,Postcode,Phone,Credit Limit,Birthday")
	expected := "<tr><td><pre>Name</pre></td><td><pre>Address</pre></td><td><pre>Postcode</pre></td><td><pre>Phone</pre></td><td align='right'><pre>Credit Limit</pre></td><td><pre>Birthday</pre></td></tr>"
	if !strings.EqualFold(result, expected) {
		t.Errorf("CSV header generated is not correct, got: %s, want: %s.", result, expected)
	}

	// table csv body tests
	result, _ = getCSVToTableRow(regexBody, "\"Johnson, John\",Voorstraat 32,3122gg,020 3849381,10000,01/01/1987")
	expected = "<tr><td><pre>Johnson, John</pre></td><td><pre>Voorstraat 32</pre></td><td><pre>3122gg</pre></td><td><pre>020 3849381</pre></td><td align='right'><pre>10000</pre></td><td><pre>01/01/1987</pre></td></tr>"
	if !strings.EqualFold(result, expected) {
		t.Errorf("CSV header generated is not correct, got: %s, want: %s.", result, expected)
	}

	// table csv body tests with challenging encoding
	result, _ = getCSVToTableRow(regexBody, "\"Smith, John\",Børkestraße 32,87823,+44 728 889838,9898.3,20/09/1999")
	expected = "<tr><td><pre>Smith, John</pre></td><td><pre>Børkestraße 32</pre></td><td><pre>87823</pre></td><td><pre>+44 728 889838</pre></td><td align='right'><pre>9898.3</pre></td><td><pre>20/09/1999</pre></td></tr>"
	if !strings.EqualFold(result, expected) {
		t.Errorf("CSV header generated is not correct, got: %s, want: %s.", result, expected)
	}
}

/**
test suite for the PRN file types to HTML converter
*/
func TestPRN(t *testing.T) {
	headerColumnIndexArray := getHeaderColumnIndex("Name            Address               Postcode Phone         Credit Limit Birthday")
	expectedHeaderColumnIndexArray := []int{74, 61, 47, 38, 16, 0}
	if len(headerColumnIndexArray) != len(expectedHeaderColumnIndexArray) {
		t.Errorf("Head column indexes array is not correct, got: %d, want: %d.", len(headerColumnIndexArray), len(expectedHeaderColumnIndexArray))
	}
	if headerColumnIndexArray[0] != expectedHeaderColumnIndexArray[0] {
		t.Errorf("Head column indexes array is not correct, got: %d, want: %d.", headerColumnIndexArray[0], expectedHeaderColumnIndexArray[0])
	}

	// table prn header tests
	result := getPRNToTableRow(headerColumnIndexArray, "Name            Address               Postcode Phone         Credit Limit Birthday")
	expected := "<tr><td><pre>Name            </pre></td><td><pre>Address               </pre></td><td><pre>Postcode </pre></td><td><pre>Phone         </pre></td><td><pre>Credit Limit </pre></td><td><pre>Birthday</pre></td></tr>"
	if !strings.EqualFold(result, expected) {
		t.Errorf("PRN header generated is not correct, got: %s, want: %s.", result, expected)
	}

	// table prn body tests
	result = getPRNToTableRow(headerColumnIndexArray, "Wicket, Steve   Mendelssohnstraat 54d 3423 ba  0313-398475          93400 19640603")
	expected = "<tr><td><pre>Wicket, Steve   </pre></td><td><pre>Mendelssohnstraat 54d </pre></td><td><pre>3423 ba  </pre></td><td><pre>0313-398475   </pre></td><td><pre>       93400 </pre></td><td><pre>19640603</pre></td></tr>"
	if !strings.EqualFold(result, expected) {
		t.Errorf("PRN header generated is not correct, got: %s, want: %s.", result, expected)
	}

	// table prn body tests with challenging encoding
	result = getPRNToTableRow(headerColumnIndexArray, "Smith, John     Børkestraße 32        87823    +44 728 889838      989830 19990920")
	expected = "<tr><td><pre>Smith, John     </pre></td><td><pre>Børkestraße 32        </pre></td><td><pre>87823    </pre></td><td><pre>+44 728 889838</pre></td><td><pre>      989830 </pre></td><td><pre>19990920</pre></td></tr>"
	if !strings.EqualFold(result, expected) {
		t.Errorf("PRN header generated is not correct, got: %s, want: %s.", result, expected)
	}
}

/**
test suite for the utils package methods
*/
func TestUtils(t *testing.T) {
	testString := "Børkestraße 32"
	result := utils.GetChars(testString, 1, 4)
	expected := "ørke"
	if !strings.EqualFold(result, expected) {
		t.Errorf("Custom substring (utils.Getchars) result is not correct, got: %s, want: %s.", result, expected)
	}
}
