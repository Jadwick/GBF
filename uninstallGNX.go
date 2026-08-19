package main

import "strings"
import "os"
import "path/filepath"

func uninstallGNX() {
	clearscreen()
	printFormattedln(Green, false, true, "Uninstalling GNX")
	printFormattedln(Green, false, false, "Checking current data.win")
	hash, err := getDataSHA256(datafile)
	if(err != nil || hash == "") {
		showErrorExit("Install failed: unable to read data.win")
	}
	id, ok := checksums[hash]
	if(ok == false) {
		showErrorExit("Install failed: Cannot determine install version. Are you using an unmodified data.win?")
	}
	var ver string = ""
	for k, v := range latestVersions {
		if(v == id) {
			ver = k;
			break
		}
	}
	if(ver != "") { //if ver isn't empty we are up-to-date
		if(strings.Contains(id, "original")) {
			printFormattedln(Red, false, false, "Vanilla Goblin Nest already found.")
			waitForEnterKey()
			return
		}
	}
	printFormattedln(Green, false, false, "data.win OK")
	printFormattedln(Green, false, false, "Checking backup data.win")
	var backupfile string
	var distributor string
	if(strings.Contains(id, "steam")) {
		backupfile = (latestVersions["original_steam"] + ".win")
		distributor = "steam"
	}
	if(strings.Contains(id, "itch")) {
		backupfile = (latestVersions["original_itch"] + ".win")
		distributor = "itch"
	}
	backupfilepath := filepath.Join(relDataDir, backupfile)
	ex, _ := exists(backupfilepath)
	if(ex == false) {
		showErrorExit("Missing backup data.win")
	}
	oldHash, err := getDataSHA256(backupfilepath)
	if(err != nil) {
		showErrorExit("Failed to read backup data.win")
	}
	oldId, ok := checksums[oldHash]
	if(ok == false) {
		showErrorExit("Backup data doesn't match any hash - corrupted?")
	}
	if(distributor == "steam") {
		if(oldId != latestVersions["original_steam"]) {
			showErrorExit("Backup data is out-of-date. Reinstall Goblin Nest with latest version through Steam.")
		}
	}
	if(distributor == "itch") {
		if(oldId != latestVersions["original_itch"]) {
			showErrorExit("Backup data is out-of-date. Reinstall Goblin Nest with latest version.")
		}
	}
	printFormattedln(Green, false, false, "Backup data.win OK")
	printFormattedln(Green, false, false, "Restoring original data.win")
	backupdata, fErr := os.ReadFile(backupfilepath)
	if(fErr != nil) {
		showErrorExit("Failed to read backup data.win")
	}
	fErr = os.WriteFile(datafile, backupdata, 0644)
	if(fErr != nil) {
		showErrorExit("Overwrite data.win failed.")
	}
	printFormattedln(Green, false, true, "GNX uninstalled successfully.")
	waitForEnterKey()
}
