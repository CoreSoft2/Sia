package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/NebulousLabs/Sia/modules"
)

var (
	renterCmd = &cobra.Command{
		Use:   "renter",
		Short: "Perform renter actions",
		Long:  "Upload and download files, or view a list of previously uploaded files.",
		Run:   wrap(renterstatuscmd),
	}

	renterUploadCmd = &cobra.Command{
		Use:   "upload [filename] [nickname] [pieces]",
		Short: "Upload a file",
		Long:  "Upload a file using a given nickname and split across 'pieces' hosts.",
		Run:   wrap(renteruploadcmd),
	}

	renterDownloadCmd = &cobra.Command{
		Use:   "download [nickname] [destination]",
		Short: "Download a file",
		Long:  "Download a previously-uploaded file to a specified destination.",
		Run:   wrap(renterdownloadcmd),
	}

	renterDownloadQueueCmd = &cobra.Command{
		Use:   "queue",
		Short: "View the download queue",
		Long:  "View the list of files that have been downloaded.",
		Run:   wrap(renterdownloadqueuecmd),
	}

	renterStatusCmd = &cobra.Command{
		Use:   "status",
		Short: "View a list of uploaded files",
		Long:  "View a list of files that have been uploaded to the network.",
		Run:   wrap(renterstatuscmd),
	}
)

func renteruploadcmd(source, nickname, pieces string) {
	err := callAPI(fmt.Sprintf("/renter/upload?source=%s&nickname=%s&pieces=%s", source, nickname, pieces))
	if err != nil {
		fmt.Println("Could not upload file:", err)
		return
	}
	fmt.Printf("Uploaded %s as '%s'.\n", source, nickname)
}

func renterdownloadcmd(nickname, destination string) {
	err := callAPI(fmt.Sprintf("/renter/download?nickname=%s&destination=%s", nickname, destination))
	if err != nil {
		fmt.Println("Could not download file:", err)
		return
	}
	fmt.Printf("Started downloading '%s' to %s.\n", nickname, destination)
}

// TODO: this should be defined elsewhere
type downloadInfo struct {
	Completed   bool
	Destination string
	Nickname    string
}

func renterdownloadqueuecmd() {
	var queue []downloadInfo
	err := getAPI("/renter/downloadqueue", &queue)
	if err != nil {
		fmt.Println("Could not get download queue:", err)
		return
	}
	if len(queue) == 0 {
		fmt.Println("No downloads to show.")
		return
	}
	fmt.Println("Download Queue:")
	for _, file := range queue {
		done := "Done"
		if !file.Completed {
			done = "... "
		}
		fmt.Printf("%s %s -> %s\n", done, file.Nickname, file.Destination)
	}
}

func renterstatuscmd() {
	status := new(modules.RentInfo)
	err := getAPI("/renter/status", status)
	if err != nil {
		fmt.Println("Could not get file status:", err)
		return
	}
	if len(status.Files) == 0 {
		fmt.Println("No files have been uploaded.")
		return
	}
	fmt.Println("Uploaded", len(status.Files), "files:")
	for _, file := range status.Files {
		fmt.Println("\t", file)
	}
}