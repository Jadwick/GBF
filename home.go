package main

import "os"

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
