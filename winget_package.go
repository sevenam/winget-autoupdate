package main

type WingetPackage struct {
	Id               string `json:"PackageIdentifier"`
	Name             string `json:"PackageName"`
	Version          string `json:"Version"`
	AvailableVersion string `json:"AvailableVersion"`
	Source           string `json:"Source"`
}
