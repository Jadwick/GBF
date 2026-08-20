package main

import "strings"
import "os"
import "path/filepath"

func restoreBackup() {
	clearscreen()
	printFormattedln(Green, false, true, "Restoring Backup")
	printFormattedln(Green, false, false, "Checking current data.win")
	hash, err := getDataSHA256(datafile)
	if(err != nil || hash == "") {
		printFormattedln(Red, false, false, "Unable to read data.win")
	}
	id, ok := checksums[hash]
	if(ok == false) {
		printFormattedln(Red, false, false, "Cannot determine installed version.")
	}
	if(strings.Contains(id, "original")) {
		msg := "Version " + id + " found."
		printFormattedln(Red, false, false, msg)
	}
	printFormattedln(Green, false, false, "Checking backups...")
	var backuppaths []string
	entries, err := os.ReadDir(relDataDir)
	if(err != nil) {
		showErrorExit("Could not read 'GBF_data' dir")
	}
	for _, e := range entries {
		name := e.Name()
		runestring := []rune(name)
	    last3 := string(runestring[len(runestring)-3:])
		if(last3 == "win") {
			fp := filepath.Join(relDataDir, name)
			hash, _ := getDataSHA256(fp)
			_, ok := checksums[hash]
			if(ok == true) {
				backuppaths = append(backuppaths, fp)
			}
	    }
	}
	var options []string = make([]string, len(backuppaths))
	for k, v := range backuppaths {
		options[k] = filepath.Base(v);
	}
	if(len(options) == 0) {
		printFormattedln(Red, false, false, "No backups found.")
		waitForEnterKey()
		return
	}
	printFormattedln(Purple, false, false, "Which version do you want to restore?")
	choice := getOption(options)
	backupfile := backuppaths[choice]
	printFormattedln(Green, false, false, "Restoring original data.win")
	backupdata, fErr := os.ReadFile(backupfile)
	if(fErr != nil) {
		showErrorExit("Failed to read backup data.win")
	}
	fErr = os.WriteFile(datafile, backupdata, 0644)
	if(fErr != nil) {
		showErrorExit("Overwrite data.win failed.")
	}
	printFormattedln(Green, false, true, "Restored backup successfully")
	waitForEnterKey()
}
