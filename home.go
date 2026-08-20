package main

import "os"

func home() {
	for {
		clearscreen()
		printFormattedln(Green, false, true, "   Goblin's Best Friend   ")
		printFormattedln(Purple, false, false, "What would you like to do?")
		var options []string = make([]string, 3)
		options[0] = "Install/Update GNX"
		options[1] = "Restore Backup (Uninstall GNX)"
		options[2] = "Quit"
		choice := getOption(options)
		switch choice {
			case 0:
				installGNX()
			case 1:
				restoreBackup()
			case 2:
				os.Exit(0)
		}
	}
}
