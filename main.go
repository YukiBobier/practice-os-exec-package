package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func main() {
	cmd := exec.Command("cat")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			buf := make([]byte, 100)
			n, err := stdout.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			output := string(buf[:n])

			fmt.Printf("%q\n", output)
		}
	}()

	go func() {
		for {
			_, err := fmt.Fprintf(stdin, "%d\n", rand.Intn(1000))
			if err != nil {
				log.Fatal(err)
			}

			time.Sleep(1 * time.Second)
		}
	}()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	s := <-c
	fmt.Println("Got signal:", s)
}
