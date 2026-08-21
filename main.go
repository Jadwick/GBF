package main

import "os"
import "path/filepath"

var latestVersions = make(map[string]string)
var checksums = make(map[string]string)
var urls = make(map[string]string)
var compatable = make(map[string]string)
var config = make(map[string]string)
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

	latestVersions["original_itch"] = "original_itch_1.36"
	latestVersions["original_steam"] = "original_steam_1.36b"
	latestVersions["patched_itch"] = "patched_itch_1.3.11.5"
	latestVersions["patched_steam"] = "patched_steam_1.3.11.5"

	checksums["b8dfaf992a2cd34b10f07bf3237961067bd315d7b1e6ff30dd8e1b003b953c58"] = "original_itch_1.33"
	checksums["919a78483e6eea1eb0e008043ae0cd7cb9fca248a2d6a9a2323139aeb2c7ffde"] = "original_steam_1.33"
	checksums["08709e8a8a474de68f2246b283945849c9d121abc7d5e2def2e7632efb235172"] = "original_itch_1.35"
	checksums["8005ebd4a226a9955da9b0a7b15cf96df911f4609beef3423126dd5a1bfcb1af"] = "original_steam_1.35"
	checksums["e5a5996c9aebb712a76bffa4199fc3cd1897005f10a4bd4dc601984a9d434dbe"] = "original_itch_1.36"
	checksums["12a0a9f92db6c9aeaf18f19599107b0f4e15c0862593852e60f778dc4a08f9ed"] = "original_steam_1.36"
	checksums["43da89d26474a2b938c1a92471d16f2b0fe0f0b12e229800ed62d68fc77f6421"] = "original_steam_1.36b"

	checksums["9a0771e11b3595317c75a89197e8d255615c1b4dff2bda0c561937e565a17c66"] = "patched_itch_1.3.9"
	checksums["0fce99c86d958e5c374753b5d144a61138a449a4a8c765590d4964aecc9eb0dd"] = "patched_steam_1.3.9"
	checksums["9171fa5d6669c8222671a01dece91326ae5ac5dcf409fd5c68cc928ad303ae53"] = "patched_itch_1.3.10"
	checksums["9ab76c0159e2301f383263daa28838b22608a8142563cf6163bca4fec250b1a3"] = "patched_steam_1.3.10"
	checksums["e6a15dbe76a4efe9cc8b3d346add06d308de778ae56b7bd4e7fab40593fb92e3"] = "patched_itch_1.3.11"
	checksums["4b1c5e92450fe18e12377af1d4dfbbd49ac3a8b61ef4bcf9b843f89195ea0dae"] = "patched_steam_1.3.11"
	checksums["bdc10150ff63f0198389de3e25dcc5267b5d7939cfefa9b6b1081e2a0b937473"] = "patched_itch_1.3.11.1"
	checksums["4ad44cea8cfe82073aa196f2c45a3be6958a5c0b4b5ed11881b74921f1cc7104"] = "patched_steam_1.3.11.1"
	checksums["9c0ff799af181046f130076d4bae256ad4d08b485c978ec9e577e3c7ec6ab6ee"] = "patched_itch_1.3.11.2"
	checksums["a22b93fd62f952681de0b46e0ccadf6b1b9c73c842fbdeec2e49432e3cf77e6c"] = "patched_steam_1.3.11.2"
	checksums["77580db641c69e88f9e0828cd4883babcf5010756e4ffb99c09cb061493e5993"] = "patched_itch_1.3.11.3"
	checksums["4a9e0741d104f6884d04868c46f34fb0331be2a178535a8ee36277e8ac4a5f88"] = "patched_steam_1.3.11.3"
	checksums["1bd74b92ec6978ef20b2a83f1dc53193636558f7c06dbbe7e92c8f8d72d211d7"] = "patched_itch_1.3.11.4"
	checksums["89c6431d6a9977606c47a319dd9af8c22ce02ce0d9fb68db33f9cd304e45d097"] = "patched_steam_1.3.11.4"
	checksums["66839d9a68d170f19136cfcd51a7d846fca3761c7a12ac49a3260ab73f7d3658"] = "patched_itch_1.3.11.5"
	checksums["3241e5a45d611ca2663f17b253506d18aa425c28353fe67d6b42481d8c31a67f"] = "patched_steam_1.3.11.5"

	urls["itch_1.3.11"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/fossil/gnx1.3.11_gn1.33_itch.fossil"
	urls["steam_1.3.11"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/fossil/gnx1.3.11_gn1.33_steam.fossil"
	urls["itch_1.3.11.1"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/fossil/gnx1.3.11.1_gn1.35_itch.fossil"
	urls["steam_1.3.11.1"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/fossil/gnx1.3.11.1_gn1.35_steam.fossil"
	urls["itch_1.3.11.2"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/fossil/gnx1.3.11.2_gn1.36_itch.fossil"
	urls["steam_1.3.11.2"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/fossil/gnx1.3.11.2_gn1.36_steam.fossil"
	urls["itch_1.3.11.3"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/fossil/gnx1.3.11.3_gn1.36b_itch.fossil"
	urls["steam_1.3.11.3"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/fossil/gnx1.3.11.3_gn1.36b_steam.fossil"
	urls["itch_1.3.11.4"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/fossil/gnx1.3.11.4_gn1.36b_itch.fossil"
	urls["steam_1.3.11.4"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/fossil/gnx1.3.11.4_gn1.36b_steam.fossil"
	urls["itch_1.3.11.5"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/fossil/gnx1.3.11.5_gn1.36b_itch.fossil"
	urls["steam_1.3.11.5"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/fossil/gnx1.3.11.5_gn1.36b_steam.fossil"

	compatable["original_itch_1.33"] = "itch_1.3.11"
	compatable["original_steam_1.33"] = "steam_1.3.11"
	compatable["original_itch_1.35"] = "itch_1.3.11.1"
	compatable["original_steam_1.35"] = "steam_1.3.11.1"
	compatable["original_itch_1.36"] = "itch_1.3.11.5"
	compatable["original_steam_1.36"] = "steam_1.3.11.2"
	compatable["original_steam_1.36b"] = "steam_1.3.11.5"

	config["updateurl"] = "https://github.com/Jadwick/GNX-Delta-Archive/raw/refs/heads/master/gbf/checksums.json"
	config["usecolor"] = "false"
	if(len(os.Args) > 1) {
		installMod(os.Args[1])
	}
	initialize()
	home()
}
