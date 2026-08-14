package main

import "fmt"
import "os"
import "io"
import "strings"
import "strconv"
import "regexp"
import "io/fs"
import "net/http"
import "errors"
import "path/filepath"
import "archive/zip"
import "encoding/json"
import "crypto/sha256"

func exists(path string) (bool, error){
	_, err := os.Stat(path)
	if (err == nil) {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func downloadFile(filepath string, url string) (err error) {
//TODO: More error checking
	var newfile = filepath
	var file_exists = false;
	b, exist_err := exists(filepath)
	if(b == true && exist_err == nil) {
		file_exists = true
		newfile = filepath + ".tmp"
	} else if (exist_err != nil) {
		return exist_err
	}

	// Create the file
	out, file_err := os.Create(newfile)
	if(file_err != nil) {
		os.Remove(newfile)
		return file_err;
	}
	defer out.Close()
	
	resp, url_err := http.Get(url)
	if(url_err != nil) {
		out.Close()
		os.Remove(newfile)
		return url_err
	}
	if(resp.StatusCode != 200) {
		out.Close()
		os.Remove(newfile)
		return errors.New("HTTP failed GET StatusCode: " + strconv.Itoa(resp.StatusCode))
	}
	defer resp.Body.Close()
	
	_, io_err := io.Copy(out, resp.Body)
	if (io_err != nil) {
		out.Close()
		os.Remove(newfile)
		return io_err
	}
	
	resp.Body.Close()
	out.Close()
	// everything worked, so rename the files and remove the old one
	if(file_exists == true){
		os.Remove(filepath)
		os.Rename(newfile, filepath)
	}
	return nil
}

func unzip(src, dest string) error {
    r, err := zip.OpenReader(src)
    if err != nil {
        return err
    }
    defer func() {
        if err := r.Close(); err != nil {
            panic(err)
        }
    }()

    os.MkdirAll(dest, 0755)

    // Closure to address file descriptors issue with all the deferred .Close() methods
    extractAndWriteFile := func(f *zip.File) error {
        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer func() {
            if err := rc.Close(); err != nil {
                panic(err)
            }
        }()

        path := filepath.Join(dest, f.Name)

        // Check for ZipSlip (Directory traversal)
        if !strings.HasPrefix(path, filepath.Clean(dest) + string(os.PathSeparator)) {
            return fmt.Errorf("illegal file path: %s", path)
        }

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, f.Mode())
        } else {
            os.MkdirAll(filepath.Dir(path), f.Mode())
            f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return err
            }
            defer func() {
                if err := f.Close(); err != nil {
                    panic(err)
                }
            }()

            _, err = io.Copy(f, rc)
            if err != nil {
                return err
            }
        }
        return nil
    }

    for _, f := range r.File {
        err := extractAndWriteFile(f)
        if err != nil {
            return err
        }
    }

    return nil
}

func updateChecksums() bool {
	printFormattedln(Green, false, false, "Attempting to download new checksums...")
	err := downloadFile(filepath.Join(DATADIR, CHECKSUMFILE), config["updateurl"])
	if (err != nil) {
		printFormattedln(Red, false, false, "Failed to download new checksums")
		printFormattedln(Red, false, false, err.Error())
		return true
	}
	downloadedChecksums = true
	return false
}

func parseChecksums() bool {
	printFormattedln(Green, false, false, "Attempting to load checksum file...")
	fileExists, exist_err := exists(filepath.Join(DATADIR, CHECKSUMFILE))
	if(exist_err != nil) {
		printFormattedln(Red, false, false, "Failed to read checksum file.")
		return true
	}
	if(fileExists == false) {
		printFormattedln(Red, false, false, "No checksum file.")
		return true
	}
	file, file_err := os.Open(filepath.Join(DATADIR, CHECKSUMFILE))
	if(file_err != nil) {
		printFormattedln(Red, false, false, "Failed to read checksum file.")
		return true
	}
	var jsonMap map[string]interface{}
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&jsonMap)
	if (err != nil) {
		file.Close()
		printFormattedln(Red, false, false, "Failed to read JSON.")
		printFormattedln(Purple, false, false, "Delete malformed checksum file?")
		ret := getYesNo(true)
		if(ret == true) {
			os.Remove(filepath.Join(DATADIR, CHECKSUMFILE))
		}
		return true
	}
	for key, value := range jsonMap {
		switch key {
			case "gbf_version":
				s, ok := value.(string)
				if(ok == true) {
					gbfJsonVersion = s
				}
			case "latest":
				m, ok := value.(map[string]interface{})
				if(ok == true) {
					for ver, id := range m {
						latestVersions[ver] = id.(string)
					}
				}
			case "checksums":
				m, ok := value.(map[string]interface{})
				if(ok == true) {
					for cs, id := range m {
						_, exists := checksums[cs]
						if (exists == false) {
							checksums[cs] = id.(string)
						}
					}
				}
			case "urls":
				m, ok := value.(map[string]interface{})
				if(ok == true) {
					for id, url := range m {
						_, exists := urls[id]
						if (exists == false) {
							urls[id] = url.(string)
						}
					}
				}
		}
	}
	file.Close()
	return false
}

func trimUrl(url string) string {
	//remove everything before the last slash
	re := regexp.MustCompile(`[^\/]+$`)
	trimmed := re.FindString(url)
	//remove anything after the first '?'
	re = regexp.MustCompile(`^[^\?]+`)
	trimmed = re.FindString(trimmed)
	return trimmed
}

func getDataSHA256(fp string) (string, error) {
	file, err := os.Open(fp)
	if(err != nil) {
		file.Close()
		return "", err
	}
	defer file.Close()
	
	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if(err != nil) {
		file.Close()
		return "", err
	}
	file.Close()
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func loadConfig() bool {
	printFormattedln(Green, false, false, "Attempting to load config file...")
	fileExists, exist_err := exists(filepath.Join(DATADIR, CONFIGFILE))
	if(exist_err != nil) {
		printFormattedln(Red, false, false, "Failed to read config file.")
		return true
	}
	if(fileExists == false) {
		printFormattedln(Red, false, false, "No config file. Creating...")
		createConfig()
		return true
	}
	file, file_err := os.Open(filepath.Join(DATADIR, CONFIGFILE))
	if(file_err != nil) {
		printFormattedln(Red, false, false, "Failed to read config file.")
		return true
	}
	var jsonMap map[string]interface{}
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&jsonMap)
	if (err != nil) {
		file.Close()
		printFormattedln(Red, false, false, "Failed to read JSON.")
		printFormattedln(Purple, false, false, "Delete malformed config file?")
		ret := getYesNo(false)
		if(ret == true) {
			os.Remove(filepath.Join(DATADIR, CONFIGFILE))
		}
		return true
	}
	for key, value := range jsonMap {
		switch key {
			case "updateurl":
				s, ok := value.(string)
				if(ok == true) {
					config["updateurl"] = s
				}
		}
	}
	file.Close()
	return false
}

func createConfig() {
	fp := filepath.Join(DATADIR, CONFIGFILE)
	file, file_err := os.Create(fp)
	if(file_err != nil) {
		os.Remove(fp)
		printFormattedln(Red, false, false, "Failed to create config file.")
		fmt.Println(file_err)
		return
	}
	defer file.Close()
	json, err := json.Marshal(config)
	if(err != nil) {
		file.Close()
		os.Remove(fp)
		printFormattedln(Red, false, false, "Failed to create JSON.")
		return
	}
	_, io_err := file.Write(json)
	if(io_err != nil) {
		file.Close()
		os.Remove(fp)
		printFormattedln(Red, false, false, "Failed to write config file.")
		return
	}
	file.Close()
	
}