package main

const CLEARSCREEN = "\033[H\033[2J"
const DATADIR = "gbf_data"
const MODDIR = "GNX_mods"
const CHECKSUMFILE = "checksums.json"
const CONFIGFILE = "config.json"
const GBF_VERSION = "1.1.1"

type Color int

const (
	Default Color = iota
	Red
	Green
	Purple
)
