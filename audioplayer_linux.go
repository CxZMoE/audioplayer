package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/CxZMoE/bass-go"
)

const (
	pidDirectory = "/tmp/cxzaudioplayer/"
)

var (
	// We need a handler to pass into vars to do something
	play        string
	recovr      string
	stop        string
	processName string
	noQuit      bool
	isLoop      bool
	pos         int
)

var isSet map[interface{}]bool = make(map[interface{}]bool, 8)

// initialize bass libray and load bass plugins
func initBass() {
	fmt.Println("Initializing bass library.")
	if bass.PluginLoad("/lib/libbassflac.so") <= 0 || bass.Init() <= 0 {
		fmt.Println("[ERR] failed to initialize bass library.")
	}

	fmt.Println("Setting default bass configurations.")

	fmt.Println("Success.")
}

// free bass library.
func freeBass() {
	fmt.Println("Freeing bass library.")
	if bass.Free() <= 0 {
		fmt.Println("[ERR] failed to free bass library.")
	}

	fmt.Println("Bass is freed.")
}

var c chan os.Signal

func init() {

	initBass()

	// create a music player process,you can specify a processID to it,if you don't stop with processID,the process will stop automatically after play is over.
	bindAlts(&play, "-1", "play a music by filename", "play", "p")
	bindAlts(&recovr, "-1", "recover a music by process name", "recover", "r")
	bindAlts(&stop, "-1", "stop a music by process name", "stop", "s")
	bindAlts(&noQuit, false, "set quit or not when play is over.", "noquit")
	bindAlts(&isLoop, false, "is playing looply.", "loop", "l")
	bindAlts(&pos, 0, "specify the position of music", "pos")
	bindAlts(&processName, "-1", "set the name of player.", "name", "n")

	flag.Parse()

	if play != "-1" {
		isSet[&play] = true
	}
	if recovr != "-1" {
		isSet[&recovr] = true
	}
	if stop != "-1" {
		isSet[&stop] = true
	}
	if noQuit != false {
		isSet[&noQuit] = true
	}
	if isLoop != false {
		isSet[&isLoop] = true
	}
	if pos != 0 {
		isSet[&pos] = true
	}
	if processName != "-1" {
		isSet[&processName] = true
	}
	if recovr != "-1" {
		isSet[&recovr] = true
	}
}

func getPP(processName string) PlayerPID {
	var pp PlayerPID
	pidFilePath := pidDirectory + processName + ".pid"
	if isNotExist(pidFilePath) {
		fmt.Println("pid file does not exist.")
		return pp
	}

	ppData, _ := ioutil.ReadFile(pidFilePath)
	json.Unmarshal(ppData, &pp)
	return pp
}

func checkIsTheOnlySettedVar(vs ...interface{}) bool {
	setCount := 0
	for k, v := range isSet {
		for _, vsv := range vs {
			if &vsv != k {
				if v {
					return false
				}
			}
			if isSet[&vsv] {
				setCount++
			}
		}
	}
	if setCount != len(vs) {
		return false
	}
	return true
}

func main() {
	// PLAY

	if play != "-1" {
		if processName == "-1" {
			processName = strconv.Itoa(int(time.Now().Unix()))
		}

		// make a signal channel
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGKILL)

		if isLoop {
			startMusicHandler(play, 1)
		} else {
			startMusicHandler(play, 0)
		}

		for {

			a := <-c

			// break when SIGINT and SIGKILL
			if a == syscall.SIGINT || a == syscall.SIGKILL {
				break
			}
		}
		freeBass()

	}

	// STOP
	if stop != "-1" {

		pidFilePath := pidDirectory + stop + ".pid"
		if isNotExist(pidFilePath) {
			fmt.Println("pid file does not exist.")
			return
		}

		pp := getPP(stop)

		err := syscall.Kill(pp.PID, syscall.SIGKILL)
		if err != nil {
			fmt.Println("[ERR] failed to stop player:", err)
		}

		os.Remove(pidFilePath)
	}

	// RECOVER
	if recovr != "-1" {
		pidFilePath := pidDirectory + recovr + ".pid"

		if isNotExist(pidFilePath) {
			fmt.Println("pid file does not exist.")
			return
		}

		pp := getPP(recovr)
		err := syscall.Kill(pp.PID, syscall.SIGKILL)
		if err != nil {
			//fmt.Println("[ERR] failed to kill existing player:", err)
		}

		fmt.Println(pp.File)
		processName = recovr
		if isLoop {
			startMusicHandler(pp.File, 1)
		} else {
			startMusicHandler(pp.File, 0)
		}

		for {

			a := <-c

			// break when SIGINT and SIGKILL
			if a == syscall.SIGINT || a == syscall.SIGKILL {
				break
			}
		}
		freeBass()

	}

}

// bind multiple name for one variable
func bindAlts(obj interface{}, value interface{}, usage string, name ...string) {
	for _, n := range name {

		switch (obj).(type) {
		case *int:
			flag.IntVar((obj).(*int), n, value.(int), usage)
			break
		case *string:
			flag.StringVar((obj).(*string), n, value.(string), usage)
			break
		case *bool:
			flag.BoolVar((obj).(*bool), n, value.(bool), usage)
		default:
			break
		}
	}
}

// PlayerPID a struct to store information about a audioplayer process.
type PlayerPID struct {
	PID    int    `json:"pid"`
	File   string `json:"file"`
	Length int    `json:"length"`
	Pos    int    `json:"pos"`
}

// pid lock
func savePID(processName string, pid int, filename string) {
	if isNotExist(pidDirectory) {
		os.MkdirAll(pidDirectory, 0755)

	}

	pp := PlayerPID{
		PID:    pid,
		File:   filename,
		Length: 0,
		Pos:    0,
	}

	ppdata, _ := json.Marshal(pp)

	if err := ioutil.WriteFile(pidDirectory+processName+".pid", ppdata, 0755); err != nil {
		fmt.Println("[ERR] failed to create pid file:", err)
	}
}

// update informations for pp
func updatePIDInfo(pidFilePath string, pp PlayerPID, length, pos int) PlayerPID {
	pp.Pos = pos
	pp.Length = length

	return pp
}

// get pid
func getPID() int {
	return os.Getpid()
}

func exit() {
	bass.Free()
	os.Exit(0)

}

// startMusicHandler start music and returns a handler of string
func startMusicHandler(path string, loop int) uint {
	// path of pid file
	pidFilePath := pidDirectory + processName + ".pid"

	if isNotExist(path) {
		fmt.Println("[ERR] file does not exist.")
		exit()
	}

	// Exist
	if !isNotExist(pidFilePath) {
		fmt.Println(pidFilePath)
		stopProcess(processName)
	}

	// create a play channel
	handle := bass.StreamCreateFile(0, path, 0, 0)
	posInt := 0
	if isSet[&recovr] {
		pp := getPP(processName)
		posInt = bass.ChannelSeconds2Bytes(handle, pp.Pos)
	} else {
		posInt = bass.ChannelSeconds2Bytes(handle, pos)
	}

	bass.ChannelSetPosition(handle, posInt, bass.BASS_POS_BYTE)

	if handle <= 0 {
		fmt.Println("[ERR] failed to create stream.")
		exit()
	}

	fmt.Println("create handler:", handle)

	bass.ChannelPlay(handle, loop)

	// create pid file so that father will know what pid to kill.
	pid := os.Getpid()
	savePID(processName, pid, path)
	fmt.Println("pid:", pid)

	// subprocess for checking play over
	go func() {
		for {

			// get basic imformations
			length := bass.ChannelGetLength(handle, bass.BASS_POS_BYTE)
			pos := bass.ChannelGetPosition(handle, bass.BASS_POS_BYTE)
			lengthSecond := bass.ChannelBytes2Seconds(handle, length)
			posSecond := bass.ChannelBytes2Seconds(handle, pos)
			fmt.Printf("length: %d pos: %d\r", length, pos)

			// We can quit
			if !noQuit {
				// exit when play is over
				if lengthSecond <= posSecond {
					os.Remove(pidFilePath)
					exit()
				}
			}

			if isNotExist(pidFilePath) {
				fmt.Println("[ERR] pid file not found.")
				exit()
			}
			pidData, _ := ioutil.ReadFile(pidFilePath)
			var pp PlayerPID
			json.Unmarshal(pidData, &pp)

			// update pp
			pp = updatePIDInfo(pidFilePath, pp, lengthSecond, posSecond)

			ppData, _ := json.Marshal(pp)
			ioutil.WriteFile(pidFilePath, ppData, 0755)

			time.Sleep(time.Millisecond * 500)
		}
	}()
	return handle
}

func stopProcess(processName string) {

	pidFilePath := pidDirectory + processName + ".pid"
	if isNotExist(pidFilePath) {
		fmt.Println("pid file does not exist.")
		return
	}

	pp := getPP(processName)
	fmt.Println("kill:", pp.PID)
	err := syscall.Kill(pp.PID, syscall.SIGKILL)
	if err != nil {
		fmt.Println("[ERR] failed to stop player:", err)
	}
}

// Check file existance
func isNotExist(path string) bool {
	_, err := os.Lstat(path)

	if os.IsNotExist(err) {
		return true
	}
	return false
}
