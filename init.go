package main

import "os"

func initialize() {
	var hadErrors = false;
	clearscreen()
	printFormattedln(Green, false, true, "Initializing GBF")

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
	printFormattedln(Green, false, true, "GBF initialized successfully. Launching...")
	if(hadErrors == true) {
		waitForEnterKey()
	}
}
