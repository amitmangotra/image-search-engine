# Search Engine for Images

## Installation Instructions:

* Make sure Go is installed, GOPATH is configured and Go projects Workspace is set up in your system.
* Create an Account on Clarifai.com to get an API key. (Recommended)
* Unzip the clarifai zip folder into `src` folder of your Go projects workspace.
* `cd clarifai`
* Put the Clarifai API key at its respective place in `main.go` and save the changes. (Recommended)
* Run `go build`.
* Run `./clarifai` and wait for the Search Engine to be ready, while the system builds an *_INVERTED INDEX_*.
* Interact with the interface on Terminal to search for images (i.e. type any term to search for images related to it), e.g. search for "exercise", "education", "architecture", "food" etc.
* Press any other alphabet except "Y" and return, to exit out of interface/system.