package main

import "strings"
import "regexp"
import "os"
import "fmt"
import "path/filepath"
import "github.com/shadowspore/fossil-delta"

func updateGNX() {
	clearscreen()
	printFormattedln(Green, false, true, "Updating GNX")
	if(downloadedChecksums == false) {
		printFormattedln(Red, false, false, "Checksums were not updated, continuing with latest known version.")
	}
	printFormattedln(Green, false, false, "Checking current data.win")
	hash, err := getDataSHA256(datafile)
	if(err != nil || hash == "") {
		showErrorExit("Install failed: unable to read data.win")
	}
	id, ok := checksums[hash]
	if(ok == false) {
		showErrorExit("Install failed: Cannot determine install version. Are you using an unmodified data.win?\nOr was there an update recently?")
	}
	var ver string = ""
	for k, v := range latestVersions {
		if(v == id) {
			ver = k;
			break
		}
	}
	if(ver != "") { //if ver isn't empty we are up-to-date
		if(strings.Contains(id, "patched")) {
			printFormattedln(Green, false, false, "GNX already up-to-date.")
			waitForEnterKey()
			return
		} else if(strings.Contains(id, "original")) {
			printFormattedln(Red, false, false, "Vanilla Goblin Nest found. Switching to install.")
			waitForEnterKey()
			installGNX()
			return
		}
	}
	if(ver == "") { //if ver is empty, something needs updating
		if(strings.Contains(id, "original")) {
			printFormattedln(Red, false, false, "Update failed: Goblin Nest version mismatch.")
			printFormattedln(Red, false, false, fmt.Sprintf("%s%s", "Have: ", id))
			printFormattedln(Red, false, false, fmt.Sprintf("%s%s%s%s", "Need: ", latestVersions["original_itch"], " or ", latestVersions["original_steam"]))
			showErrorExit("Upgrade or downgrade to required version to install GNX (Through Steam or Itch.io)")
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
	var urlID string
	if(strings.Contains(id, "steam")) {
		re := regexp.MustCompile(`[^_]+$`)
		gnxVer := re.FindString(latestVersions["patched_steam"])
		urlID = "steam_" + gnxVer
	} else if(strings.Contains(id, "itch")) {
		re := regexp.MustCompile(`[^_]+$`)
		gnxVer := re.FindString(latestVersions["patched_itch"])
		urlID = "itch_" + gnxVer
	}
	var url string
	url, ok = urls[urlID]
	if(ok == false) {
		showErrorExit("URL to patch file wasn't located.")
	}
	printFormattedln(Green, false, false, "Downloading delta file...")
	filename := trimUrl(url)
	if(filename == "") {
		showErrorExit("Malformed URL in JSON.")
	}
	dlErr := downloadFile(filepath.Join(relDataDir, filename), url)
	if(dlErr != nil) {
		showErrorExit("Failed to download delta file.")
	}
	printFormattedln(Green, false, false, filename + " downloaded.")
	printFormattedln(Green, false, false, "Attempting delta patch...")

	patchfile := filepath.Join(relDataDir, filename)
	patchdata, fErr := os.ReadFile(patchfile)
	if(fErr != nil) {
		showErrorExit("Cannot open delta file.")
	}

	datadata, fErr := os.ReadFile(backupfilepath)
	if(fErr != nil) {
		showErrorExit("Cannot open backup data.win")
	}

	applied, fErr := fdelta.Apply(datadata, patchdata)
	if(fErr != nil) {
		showErrorExit("Failed to apply delta patch.")
	}

	fErr = os.WriteFile(datafile, applied, 0644)
	if(fErr != nil) {
		showErrorExit("Failed to overwrite data.win.")
	}
	doublecheck, err := getDataSHA256(datafile)
	if(err != nil) {
		showErrorExit("data.win failed doublecheck.")
	}
	dcID := checksums[doublecheck]
	if(dcID == "") {
		showErrorExit("Patched data.win does not match hash.")
	}
	if !(dcID == latestVersions["patched_itch"] || dcID == latestVersions["patched_steam"]) {
		printFormattedln(Red, false, false, "Patch did something weird.")
	}
	printFormattedln(Green, false, false, "Delta patch successful.")
	printFormattedln(Green, false, false, "Checking for mod directory...")
	b, err := exists(MODDIR)
	if(b == true) {
		printFormattedln(Green, false, false, "Mod directory located")
	} else if (b == false && err == nil) {
		printFormattedln(Red, false, false, "Mod directory missing, creating...")
		err := os.Mkdir(relModDir, 0755)
		if (err != nil && !os.IsExist(err)) {
			showErrorExit("Permission Error.")
		}
		printFormattedln(Green, false, false, "Mod directory created successfully")
	} else {
		printFormattedln(Red, false, true, "Unable to make 'GNX_mods' directory.")
	}
	printFormattedln(Green, false, true, "GNX installed successfully.")
	printFormattedln(Purple, false, false, "Unzip mods into 'GNX_mods'.")
	waitForEnterKey()
}
