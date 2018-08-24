/*
 * Version 2018-08-24. Initial version: 2015-11-03
 * Simpel stopwatch 
 * usage: 
 * run stopper
 * space - lap time
 * q - quit program
 * 
 * Currently the program works properly only under Linux.
 * On BSD (perhaps also on the macOS) keyboard reading (space, q) doesn't work,
 * the problem seems to in go channels. 
 */

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"
)

func formattedTime(t int64) string {
	s := ""
	//t = deciseconds
	days := t / 864000
	hours := (t - days*864000) / 36000
	minutes := (t - days*864000 - hours*36000) / 600
	seconds := (t - days*864000 - hours*36000 - minutes*600) / 10
	ds := t - days*864000 - hours*36000 - minutes*600 - seconds*10
	if days > 0 {
		s = strconv.FormatInt(days, 10) + "d"
	}
	if hours > 0 || days > 0 {
		s = s + strconv.FormatInt(hours, 10) + "h"
	}
	if minutes > 0 || days > 0 || hours > 0 {
		s = s + strconv.FormatInt(minutes, 10) + "m"
	}
	s = s + strconv.FormatInt(seconds, 10) + "." + strconv.FormatInt(ds, 10) + "s"
	return s
}

func terminalBackToNormal() {
	//Back to normal terminal
	exec.Command("stty", "-F", "/dev/tty", "-raw").Run()
	//Toores ja vana
	//exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	//Muudatus 2018-08-24 22:05
        exec.Command("stty", "-F", "/dev/tty", "echo", "-brkint", "-ignpar", "-istrip").Run()
}

func terminalToRaw() {
	//Disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	//Do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

func main() {
	terminalToRaw()

	/* Capture ctrl C
	http://stackoverflow.com/questions/11268943/golang-is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in
	*/
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			//sig is a ^C, handle it here
			//Terminal back to normal
			terminalBackToNormal()
			fmt.Println(sig)
			os.Exit(0)
		}
	}()
	//End of ctrl C capture :)

	start_time := time.Now().UnixNano()
	var total, last int64 = 0, 0

	/*
		http://stackoverflow.com/questions/15159118/read-a-character-from-standard-input-in-go-without-pressing-enter
	*/
	go func(int64, int64) {
		var b []byte = make([]byte, 1)
		count := 1
		for {
			os.Stdin.Read(b)
			if string(b) == " " {
				fmt.Println(count, "Total time: ", formattedTime(total),
					" Lap: ", formattedTime(total-last))
				count++
				last = total
			}
			if string(b) == "q" {
				terminalBackToNormal()
				fmt.Println()
				os.Exit(0)
			}
		}
	}(total, last)

	for {
		time.Sleep(100 * time.Millisecond)
		total = (time.Now().UnixNano() - start_time) / 1e8
		fmt.Print("                                          \r")
		fmt.Print("Runtime: ", formattedTime(total), " (Lap: ",
			formattedTime(total-last), ")\r")
	}
}
