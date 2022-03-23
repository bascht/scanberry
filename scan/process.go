package scan

import (
	"fmt"
	"path/filepath"
	"strings"
)

func Process(basedir string, document *Document) {

	// time.Sleep(5 * time.Second)
	fmt.Println("STARTING SCAN OF "+document.Name+" with %v to "+basedir, document.Args())

	var feed string

	if document.Duplex {
		feed = "ADF Duplex"
	} else {
		feed = "ADF Front"
	}

	var commands []Command

	commands = append(commands, Command{Name: "scanimage", Args: []string{"--device-name=epjitsu", "--format=tiff", "--batch=" + filepath.Join(basedir, document.FullName()+"-%d.tif"), "--source=" + feed, "--mode=Gray", "--resolution=300"}})
	commands = append(commands, Command{Name: "convert", Args: []string{"-verbose", filepath.Join(basedir, document.FullName() + "*tif"), filepath.Join(basedir, document.FullName()+"_input.pdf")}})
	commands = append(commands, Command{Name: "convert", Args: []string{"-verbose", filepath.Join(basedir, document.FullName() + "-1.tif"), filepath.Join(basedir, "downloads", document.FullName()+".thumbnail.jpg")}})
	commands = append(commands, Command{Name: "ocrmypdf", Args: []string{"--language=deu", filepath.Join(basedir, document.FullName()+"_input.pdf"), filepath.Join(basedir, "downloads", document.FullName()+".pdf")}})
	commands = append(commands, Command{Name: "cp", Args: []string{filepath.Join(basedir, "downloads", document.FullName()+".pdf"), "/mnt/himbeerkompott/home/bascht/Documents/Scans"}})

		<- document.Events
	for _, command := range commands {
		scanner := command.GetScanner()
		command.Start()

		go func() {

			document.Events <- Event{Message: "Starting with args " + strings.Join(command.Args, " "), Type: command.Name }
			for scanner.Scan() {
				document.Events <- Event{Message: scanner.Text(), Type: command.Name }
			}

		}()
		command.Wait()
		if command.cmd.ProcessState.ExitCode() != 0 {
			document.Events <- Event{Message: "The command did not run successfully :-(", Type: "error" }
			close(document.Events)
			return
		}
	}


	document.Events <- Event{Message: "Scan finished successfully", Type: "done" }
	// // Open original file
	// original, err := os.Open(filepath.Join(basedir, document.FullNameWithExtension()))
	// if err != nil {
	// 	document.Events <- Event{Message: "Alles kaputt: " + err.Error(), Type: "error"}
	// 	// log.Fatal(err)
	// }
	// defer original.Close()

	// // Create new file
	// newPath, err := filepath.Abs("./downloads/" + document.FullNameWithExtension())
	// new, err := os.Create(newPath)
	// if err != nil {
	// 	document.Events <- Event{Message: "Alles kaputt", Type: "error"}
	// 	// log.Fatal(err)
	// }
	// defer new.Close()

	// //This will copy
	// bytesWritten, err := io.Copy(new, original)
	// if err != nil {
	// 	document.Events <- Event{Message: "Alles kaputt", Type: "error"}
	// 	// log.Fatal(err)
	// }
	// fmt.Printf("Bytes Written: %d\n", bytesWritten)

	close(document.Events)
}
