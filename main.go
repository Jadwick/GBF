package main

import "os"
import "path/filepath"

var latestVersions = make(map[string]string)
var checksums = make(map[string]string)
var urls = make(map[string]string)
var config = make(map[string]string)
var useColor bool = true
var downloadedChecksums = false
var gbfJsonVersion string = ""
var relExeDir string
var relDataDir string
var relModDir string
var datafile string

func main() {
	exe, err := os.Executable()
	if(err != nil) {
		showErrorExit(err.Error())
	}
	relExeDir = filepath.Dir(exe)
	relDataDir = filepath.Join(relExeDir, DATADIR)
	relModDir = filepath.Join(relExeDir, MODDIR)
	datafile = filepath.Join(relExeDir, "data.win")
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

	config["updateurl"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/gbf/checksums.json"
	if(len(os.Args) > 1) {
		installMod(os.Args[1])
	}
	initialize()
	home()
}
