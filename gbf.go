package main

import "fmt"
import "os"
import "io"
import "strings"
import "regexp"
import "path/filepath"
import "github.com/shadowspore/fossil-delta"

type Color int

const (
	Default Color = iota
	Red
	Green
	Purple
)

const CLEARSCREEN = "\033[H\033[2J"
const DATADIR = "gbf_data"
const MODDIR = "GNX_mods"
const CHECKSUMFILE = "checksums.json"
const CONFIGFILE = "config.json"
const GBF_VERSION = "1.0.0"

var latestVersions = make(map[string]string)
var checksums = make(map[string]string)
var urls = make(map[string]string)
var config = make(map[string]string)
var useColor bool = true
var downloadedChecksums = false
var gbfJsonVersion string = ""

var datafile string = "data.win"

func main() {
	useColor = doesTerminalSupportColor()
	
	latestVersions["original_itch"] = "original_itch_1.33"
	latestVersions["original_steam"] = "original_steam_1.33"
	latestVersions["patched_itch"] = "patched_itch_1.3.11"
	latestVersions["patched_steam"] = "patched_steam_1.3.11"
	
	checksums["b8dfaf992a2cd34b10f07bf3237961067bd315d7b1e6ff30dd8e1b003b953c58"] = "original_itch_1.33"
	checksums["919a78483e6eea1eb0e008043ae0cd7cb9fca248a2d6a9a2323139aeb2c7ffde"] = "original_steam_1.33"
	checksums["9a0771e11b3595317c75a89197e8d255615c1b4dff2bda0c561937e565a17c66"] = "patched_itch_1.3.9"
	checksums["0fce99c86d958e5c374753b5d144a61138a449a4a8c765590d4964aecc9eb0dd"] = "patched_steam_1.3.9"
	checksums["9171fa5d6669c8222671a01dece91326ae5ac5dcf409fd5c68cc928ad303ae53"] = "patched_itch_1.3.10"
	checksums["9ab76c0159e2301f383263daa28838b22608a8142563cf6163bca4fec250b1a3"] = "patched_steam_1.3.10"
	checksums["e6a15dbe76a4efe9cc8b3d346add06d308de778ae56b7bd4e7fab40593fb92e3"] = "patched_itch_1.3.11"
	checksums["4b1c5e92450fe18e12377af1d4dfbbd49ac3a8b61ef4bcf9b843f89195ea0dae"] = "patched_steam_1.3.11"
	
	urls["itch_1.3.11"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/fossil/gnx1.3.11_gn1.33_itch.fossil"
	urls["steam_1.3.11"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/fossil/gnx1.3.11_gn1.33_steam.fossil"
	
	config["updateurl"] = "https://raw.githubusercontent.com/Jadwick/GNX-Delta-Archive/refs/heads/master/gbf/checksums.json"
	initialize()
	home()
}

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
	b, err = exists(DATADIR)
	if(b == true) {
		printFormattedln(Green, false, false, "Data directory located")
	} else if (b == false && err == nil) {
		hadErrors = true
		printFormattedln(Red, false, false, "Data directory missing, creating...")
		err := os.Mkdir(DATADIR, 0755)
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
	}
	printFormattedln(Green, false, true, "GBF initialized successfully. Launching...")
	if(hadErrors == true) {
		waitForEnterKey()
	}
}

func home() {
	for {
		clearscreen()
		printFormattedln(Green, false, true, "   Goblin's Best Friend   ")
		printFormattedln(Purple, false, false, "What would you like to do?")
		var options []string = make([]string, 4)
		options[0] = "Install GNX"
		options[1] = "Update GNX"
		options[2] = "Uninstall GNX"
		options[3] = "Quit"
		choice := getOption(options)
		switch choice {
			case 1:
				installGNX()
			case 2:
				updateGNX()
			case 3:
				uninstallGNX()
			case 4:
				os.Exit(0)
			
		}
	}
}

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
		showErrorExit("Install failed: Cannot determine install version. Are you using an unmodified data.win?")
	}
	var ver string = ""
	for k, v := range latestVersions {
		if(v == id) {
			ver = k;
			break
		}
	}
	if(ver == "") { //if ver is empty, something needs updating
		if(strings.Contains(id, "original")) {
			showErrorExit("Install failed: Goblin Nest out-of-date. Please update to latest version. (Through Steam or Itch.io)")
		} else {
			printFormattedln(Red, false, false, "Old GNX version found. Switching to update.")
				waitForEnterKey()
				updateGNX()
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
			printFormattedln(Green, false, false, "GNX already installed.")
			printFormattedln(Red, false, false, "Checksums were unable to be retrieved to verify latest version.")
			waitForEnterKey()
			home()
			return
		}
	}
	printFormattedln(Green, false, false, "Original data.win found, attempting GNX install...")
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
	dlErr := downloadFile(filepath.Join(DATADIR, filename), url)
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
	newfile := filepath.Join(DATADIR, id + ".win")
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
	
	patchfile := filepath.Join(DATADIR, filename)
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
		err := os.Mkdir(MODDIR, 0755)
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
			showErrorExit("Install failed: Goblin Nest out-of-date. Please update to latest version. (Through Steam or Itch.io)")
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
	backupfilepath := filepath.Join(DATADIR, backupfile)
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
	dlErr := downloadFile(filepath.Join(DATADIR, filename), url)
	if(dlErr != nil) {
		showErrorExit("Failed to download delta file.")
	}
	printFormattedln(Green, false, false, filename + " downloaded.")
	printFormattedln(Green, false, false, "Attempting delta patch...")
	
	patchfile := filepath.Join(DATADIR, filename)
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
		err := os.Mkdir(MODDIR, 0755)
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
	backupfilepath := filepath.Join(DATADIR, backupfile)
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


func showChecksums() {
	for key, value := range checksums {
		fmt.Println(key, value)
	}
	fmt.Scanln()
}