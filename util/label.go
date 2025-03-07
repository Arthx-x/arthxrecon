package util

import "github.com/rs/zerolog"

// Labels centralizes all application messages and texts.
var (

	// MESSAGENS
	AppName          = "ArthxRecon"
	WelcomeMessage   = "Welcome to ArthxRecon!"
	AppDescription   = "A modular recon tool for pentesting"
	CmdUsage         = "Use this tool with the provided subcommands to run various enumeration modules."
	ErrInvalidIP     = "Invalid IP address provided."
	LogFileNotFound  = "Log file not found, using console output."
	ErrInvalidTarget = "No valid target provided. Use --target to specify IP(s), CIDR, or a file containing targets."

	FatalErrHD         = "Host Discovery Failed!"
	FatalErrPS         = "Port Scan Failed!"
	FallbackConsoleMsg = "Failed to open log file, using console output" // FallbackConsoleMsg is the message used when the log file cannot be opened.
	HDAppDescription   = "Executes host discovery using Nmap"

	//CONST
	DefaultTimeFormat     = zerolog.TimeFormatUnix // DefaultTimeFormat defines the default time field format for Zerolog.
	ConfigFilePath        = "config/config.toml"   // ConfigFilePath is the path to the configuration file.
	HostDiscoveryName     = "hostDiscovery"
	PortScanName          = "portScan"
	HostDiscoveryFlagNmap = "-PS22,2222,53,80,443,445,3389"
)
