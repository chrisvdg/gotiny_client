# GoTiny client

This repository contains a simple client for the [GoTiny](https://github.com/chrisvdg/gotiny) API

## Usage

```go
package main

import (
	"fmt"
	"log"

	gotinyclient "github.com/chrisvdg/gotiny_client"
	"github.com/pkg/errors"
)

func main() {
	baseURL := "http://127.0.0.1:8080"
	token := "AVerySecretToken"

	cl, err := gotinyclient.New(baseURL, token, token)
	if err != nil {
		log.Fatal(err)
	}

	id := "google"
	url := "google.com"
	url2 := "duckduckgo.com"
	url3 := "http://qwant.com"

	// Create a new entry with ID and URL
	newEntry, err := cl.CreateEntry(id, url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created entry ", newEntry.ID, newEntry.URL)

	// Create a new entry without ID (generated)
	anotherEntry, err := cl.CreateEntry("", url2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created entry ", anotherEntry.ID, anotherEntry.URL)

	// Update entry
	err = cl.UpdateEntry(anotherEntry.ID, url3)
	if err != nil {
		log.Fatal(err)
	}

	// List current entries
	currentEntries, err := cl.ListEntries()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listing current entries")
	for _, e := range currentEntries {
		fmt.Println(e.ID, e.URL)
	}
	fmt.Println("Listed current entries")

	// delete both entries
	for _, id := range []string{newEntry.ID, anotherEntry.ID} {
		err := cl.DeleteEntry(id)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Try to fetch one of the deleted entries
	_, err = cl.GetEntry(newEntry.ID)
	if errors.Cause(err) == gotinyclient.ErrEntryNotFound {
		fmt.Println("Entry was succesfully deleted")
	} else if err == nil {
		fmt.Println("Entry still exists")
	} else {
		log.Fatal(err)
	}
}

```
