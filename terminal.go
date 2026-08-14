package main

import "fmt"
import "os"
import "strings"
import "strconv"

func clearscreen() {
	if(useColor == true) {
		fmt.Print(CLEARSCREEN)
	} else {
		systemClearscreen()
	}
}

func waitForEnterKey() {
	fmt.Println("Press Enter to continue...")
	fmt.Scanln()
}

func showErrorExit(error string) {
	printFormattedln(Red, false, false, error)
	printFormattedln(Red, false, false, "Exiting.")
	waitForEnterKey()
	os.Exit(0)
}

func printFormattedln(color Color, underline bool, invert bool, text string) {
	if(useColor == true) {
		colorcode := ""
		if(underline == true) {
			colorcode += "[4m"
		}
		switch(color) {
			case Red:
				colorcode += "[91m"
			case Green:
				colorcode += "[32m"
			case Purple:
				colorcode += "[95m"
			default:
				colorcode += ""
		}
		if(invert == true) {
			colorcode += "[7m"
		}
		fmt.Println(colorcode + text + "[0m")
	} else {
		fmt.Println(text)
	}
	
}

func getOption(options []string) int {
	var i = 1;
	for _, option := range options {
		printFormattedln(Default, false, false, strconv.Itoa(i) + ") " + option)
		i += 1;
	}
	var input string
	fmt.Scanln(&input)
	var digit int = 0
	digit, err := strconv.Atoi(input)
	if (err != nil) {
		printFormattedln(Red, false, false, "Invalid, choose again.")
		return getOption(options)
	}
	if (digit > 0 && digit < i) {
		return digit
	} else {
		printFormattedln(Red, false, false, "Invalid, choose again.")
		return getOption(options)
	}
}

// returns true on yes, false on no, 'def' is the default when user puts in a invalid option
func getYesNo(def bool) bool {
	if(def == true) {
		printFormattedln(Default, false, false, "Y/n (default yes)")
	} else {
		printFormattedln(Default, false, false, "y/N (default no)")
	}
	var input string
	fmt.Scanln(&input)
	if(strings.ToLower(input) == "y") {
		return true
	}
	if(strings.ToLower(input) == "n") {
		return false
	}
	return def
}