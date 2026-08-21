package main

import "os"
import "path/filepath"

func initialize() {
	var hadErrors = false;
	clearscreen()
	printFormattedln(Green, false, true, "   Initializing GBF   ")

	// ensure data.win exists in local directory
	b, err := exists(datafile)
	if(b == true) {
		printFormattedln(Green, false, false, "data.win located")
	} else if (b == false && err == nil) {
		hadErrors = true
		showErrorExit("data.win missing")
	} else {
		hadErrors = true
		showErrorExit("Permission Error.")
	}

	// ensure we have a 'gnxpatcher_data' directory
	b, err = exists(relDataDir)
	if(b == true) {
		printFormattedln(Green, false, false, "Data directory located")
	} else if (b == false && err == nil) {
		hadErrors = true
		printFormattedln(Red, false, false, "Data directory missing, creating...")
		err := os.Mkdir(relDataDir, 0755)
		if (err != nil && !os.IsExist(err)) {
			showErrorExit("Permission Error.")
		}
		printFormattedln(Green, false, false, "Data directory created successfully")
	} else {
		hadErrors = true
		showErrorExit("Permission Error.")
	}

	ex, err := exists(filepath.Join(relDataDir, CONFIGFILE))
	if(err != nil) {
		showErrorExit(err.Error())
	}
	if(ex == false) {
		initialSetup()
	}
	bErr := loadConfig()
	hadErrors = hadErrors || bErr
	bErr = updateChecksums()
	hadErrors = hadErrors || bErr
	bErr = parseChecksums()
	hadErrors = hadErrors || bErr
	if !(gbfJsonVersion == "" || gbfJsonVersion == GBF_VERSION) {
		hadErrors = true
		printFormattedln(Red, false, true, "New version of GBF available. ")
		printFormattedln(Purple, false, false, "https://github.com/Jadwick/GBF")
	}
	printFormattedln(Green, false, true, "   GBF initialized successfully. Launching...   ")
	if(hadErrors == true) {
		waitForEnterKey()
	}
}

func initialSetup() {
	config["usecolor"] = "true"
	printFormattedln(Green, false, true, "Running first launch setup (one-time).")
	printFormattedln(Purple, false, false, "Is the line above green, or do you only see weird numbers?")
	options := make ([]string, 2)
	options[0] = "I can see the different colors."
	options[1] = "I can only see numbers."
	choice := getOption(options)
	switch choice {
		case 0:
			config["usecolor"] = "true"
			printFormattedln(Green, false, false, "Terminal supports color.")
		case 1:
			config["usecolor"] = "false"
			printFormattedln(Red, false, false, "Terminal does not support color.")
	}
	createConfig()
}
