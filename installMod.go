package main

import "os"
import "github.com/google/uuid"
import "path/filepath"
import "io/fs"
import "strings"
import "encoding/json"

func installMod(path string) {
    runestring := []rune(path)

    last3 := string(runestring[len(runestring)-3:])

    if(last3 != "zip") {
        showErrorExit("Filetype unsupported.")
    }

    b, err := exists(datafile)
    if(b == false) {
        showErrorExit("data.win missing")
    }
    if(err != nil) {
        showErrorExit("Permission error.")
    }
    b, err = exists(relModDir)
    if(b == false) {
        showErrorExit("No 'GNX_mods' directory found. Is GNX installed?")
    }
    if(err != nil) {
        showErrorExit("Permission error.")
    }
    b, err = exists(relDataDir)
    if(b == false) {
        showErrorExit("No 'GBF_data' directory found. Is GBF in the right directory?")
    }
    if(err != nil) {
        showErrorExit("Permission error.")
    }
    printFormattedln(Green, false, false, "Attempting to install mod from:")
    printFormattedln(Green, false, false, path)
    unique := uuid.New().String()

    tempdir := filepath.Join(relDataDir, unique)
    b, err = exists(tempdir)
    if(b == true) {
        showErrorExit("Temp directory already exists!")
    }
    if(err != nil) {
        showErrorExit("Permission error.")
    }

    err = os.Mkdir(tempdir, 0755)
    if (err != nil && !os.IsExist(err)) {
        showErrorExit("Permission Error.")
    }

    unzip(path, tempdir)

    var modpath string
    var manifestroot string
    var success bool = false;
    var modname string
    var modid string

    filepath.WalkDir(tempdir, func(path string, d fs.DirEntry, err error) error {
        if(!strings.Contains(strings.ToLower(path), "manifest.json")) {
            return nil
        }
        printFormattedln(Green, false, false, "Manifest found!")
        manifestroot = filepath.Dir(path)
        //Read json
        file, err := os.Open(path)
        if(err != nil) {
            printFormattedln(Red, false, false, "Failed to read manifest.")
            return filepath.SkipAll
        }
        var jsonMap map[string]interface{}
        decoder := json.NewDecoder(file)
        err = decoder.Decode(&jsonMap)
        if (err != nil) {
            file.Close()
            printFormattedln(Red, false, false, "Failed to read JSON.")
            return filepath.SkipAll
        }
        file.Close()

        for key, value := range jsonMap {
    		switch key {
                case "mod_id":
    				s, ok := value.(string)
    				if(ok == true) {
                        modid = s
    				}
                case "name":
                    s, ok := value.(string)
                    if(ok == true) {
                        modname = s
                    }
    		}
    	}

        if(modid == ""){
            modid = modname
        }
        if(modname == "") {
            return filepath.SkipAll
        }

        modpath = filepath.Join(relModDir, modid)
        success = true
        return filepath.SkipAll
    })
    if(err != nil) {
        cleandir(tempdir)
        showErrorExit("Failed directory walk.")
    }

    if(success == false) {
        cleandir(tempdir)
        showErrorExit("No readable manifest found!")
    }

    b, err = exists(modpath)
    if(err != nil) {
        cleandir(tempdir)
        showErrorExit("Permission Error.")
    }
    if(b == true) {
        printFormattedln(Red, false, false, "Mod directory already exists, overwrite?")
        if(getYesNo(false) == false) {
            cleandir(tempdir)
            showErrorExit("Did not install mod.")
        }
        err = os.RemoveAll(modpath)
        if(err != nil) {
            cleandir(tempdir)
            showErrorExit(err.Error())
        }
    }

    msg := "Copying mod files to: " + modpath
    printFormattedln(Green, false, false, msg)

    err = os.CopyFS(modpath, os.DirFS(manifestroot))
    if(err != nil) {
        printFormattedln(Red, false, false, err.Error())
    }

    msg = modname + " installed sucessfully. Cleaning up."
    printFormattedln(Green, false, false, msg)
    cleandir(tempdir)
    waitForEnterKey()
    os.Exit(0)
}

func cleandir(path string) {
    err := os.RemoveAll(path)
    if(err != nil){
        showErrorExit(err.Error())
    }
}
