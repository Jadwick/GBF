package main

import "strings"
import "os"
import "fmt"
import "path/filepath"
import "io"
import "github.com/shadowspore/fossil-delta"

func installGNX() {
	clearscreen()
	printFormattedln(Green, false, true, "Installing GNX")
	printFormattedln(Green, false, false, "Checking data.win")
	hash, err := getDataSHA256(datafile)
	if(err != nil || hash == "") {
		showErrorExit("Install failed: unable to read data.win")
	}
	id, ok := checksums[hash]
	if(ok == false) {
		showErrorExit("Install failed: Cannot determine install version. Are you using an unmodified data.win?\nOr was there an update recently?")
	}
	msg := "Found " + id
	printFormattedln(Green, false, false, msg)
	var ver string = ""
	for k, v := range latestVersions {
		if(v == id) {
			ver = k;
			break
		}
	}
	if(ver == "") { //if ver is empty, something needs updating
		if(strings.Contains(id, "original")) {
			printFormattedln(Red, false, true, "Old Goblin Nest version found")
			printFormattedln(Red, false, false, fmt.Sprintf("%s%s", "Have: ", id))
			printFormattedln(Red, false, false, fmt.Sprintf("%s%s%s%s", "Latest: ", latestVersions["original_itch"], " or ", latestVersions["original_steam"]))
			gv, ok := compatable[id]
			if(ok == true) {
				msg := "Out-of-date GNX found for old version: " + gv
				printFormattedln(Green, false, false, msg)
				printFormattedln(Purple, false, false, "Install old version?")
				ans := getYesNo(false)
				if(ans == false) {
					printFormattedln(Red, false, false, "User canceled install.")
					return
				}
			} else {
				showErrorExit("No GNX found for this version.")
			}
		} else {
			printFormattedln(Red, false, false, "Old GNX version found.")
			printFormattedln(Purple, false, false, "Would you like to update?")
			choice := getYesNo(true)
			if(choice == true) {
				updateGNX()
			} else {
				printFormattedln(Red, false, false, "User canceled update.")
				waitForEnterKey()
			}
			return
		}
	}
	if(strings.Contains(id, "patched")) {
		if(downloadedChecksums == true) {
			printFormattedln(Green, false, false, "GNX already installed and up-to-date.")
			waitForEnterKey()
			home()
			return
		} else {
			printFormattedln(Green, false, false, "GNX already installed and assumed up-to-date.")
			printFormattedln(Red, false, false, "Checksums were unable to be retrieved to verify latest version.")
			waitForEnterKey()
			home()
			return
		}
	}
	printFormattedln(Green, false, false, "Attempting GNX install...")
	var urlID string
	var url string
	urlID, ok = compatable[id]
	if(ok == false) {
		showErrorExit("Compatable GNX version not found.")
	}
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
	printFormattedln(Green, false, false, "Backing up original data.win")
	fin, fErr := os.Open(datafile)
	if(fErr != nil) {
		showErrorExit("Failed to open data.win")
	}
	defer fin.Close()
	newfile := filepath.Join(relDataDir, id + ".win")
	fout, fErr := os.Create(newfile)
	if(fErr != nil) {
		fin.Close()
		os.Remove(newfile)
		showErrorExit("Failed to create duplicate data.win")
	}
	defer fout.Close()

	_, fErr = io.Copy(fout, fin)
	if(fErr != nil) {
		fin.Close()
		os.Remove(newfile)
		fout.Close()
		showErrorExit("Failed copying data.win")
	}
	fout.Close()
	fin.Close()
	printFormattedln(Green, false, false, "Attempting delta patch...")

	patchfile := filepath.Join(relDataDir, filename)
	patchdata, fErr := os.ReadFile(patchfile)
	if(fErr != nil) {
		showErrorExit("Cannot open delta file.")
	}

	datadata, fErr := os.ReadFile(datafile)
	if(fErr != nil) {
		showErrorExit("Cannot open data.win")
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
	dcID, ok := checksums[doublecheck]
	if(ok == false) {
		showErrorExit("Patch failed - data.win corrupted.")
	}
	msg = "GNX " + dcID + " installed."
	printFormattedln(Green, false, false, msg)
	printFormattedln(Green, false, false, "Checking for mod directory...")
	b, err := exists(relModDir)
	if(b == true) {
		printFormattedln(Green, false, false, "Mod directory located")
	} else if (b == false && err == nil) {
		printFormattedln(Red, false, false, "Mod directory missing, creating...")
		err := os.Mkdir(relModDir, 0755)
		if (err != nil && !os.IsExist(err)) {
			printFormattedln(Red, false, false, err.Error())
		}
		printFormattedln(Green, false, false, "Mod directory created successfully")
	} else {
		printFormattedln(Red, false, true, "Unable to make 'GNX_mods' directory.")
	}
	printFormattedln(Green, false, true, "GNX installed successfully")
	printFormattedln(Purple, false, false, "Drop mods onto GBF to auto-install them.")
	waitForEnterKey()
}
