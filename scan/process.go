package scan

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func Process(basedir string, document *Document) string {

	// time.Sleep(5 * time.Second)
	fmt.Println("STARTING SCAN OF "+document.Name+" with %v to "+basedir, document.Args())

	var feed string

	if document.Duplex {
		feed = "ADF Duplex"
	} else {
		feed = "ADF Front"
	}

	var commands []Command

	commands = append(commands, Command{Name: "scanimage", Args: []string{"--device-name=epjitsu:libusb:001:003", "--format=tiff", "--batch=" + filepath.Join(basedir, document.FullName()+"-%d.tif"), "--source=" + feed, "--mode=Gray", "--resolution=300"}})
	commands = append(commands, Command{Name: "convert", Args: []string{filepath.Join(basedir, document.FullName()) + "*tif", filepath.Join(basedir, document.FullName()+"_input.pdf")}})
	commands = append(commands, Command{Name: "convert", Args: []string{filepath.Join(basedir, document.FullName()) + "-1.tif", filepath.Join(basedir, document.FullName()+".thumbnail.jpg")}})
	commands = append(commands, Command{Name: "ocrmypdf", Args: []string{"--language=deu", filepath.Join(basedir, document.FullName()+"_input.pdf"), filepath.Join(basedir, document.FullName()+".pdf")}})

	for _, command := range commands {
		done := make(chan bool)
		command.Start()
		go func() {
			document.Events <- command.Name + " Startet"

			scanner := command.NewScanner()

			for scanner.Scan() {
				document.Events <- scanner.Text()
			}
			document.Events <- command.Name + " Ist fast fertig"
			done <- true
		}()
		document.Events <- command.Name + " Ist gestartet"
		command.Wait()
		<-done
	}

	// Open original file
	original, err := os.Open(filepath.Join(basedir, document.FullNameWithExtension()))
	if err != nil {
		log.Fatal(err)
	}
	defer original.Close()

	// Create new file
	newPath, err := filepath.Abs("./downloads/" + document.FullNameWithExtension())
	new, err := os.Create(newPath)
	if err != nil {
		log.Fatal(err)
	}
	defer new.Close()

	//This will copy
	bytesWritten, err := io.Copy(new, original)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bytes Written: %d\n", bytesWritten)

	close(document.Events)
	return filepath.Join(basedir, document.FullNameWithExtension())
}
