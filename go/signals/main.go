// Package signals
/*
General description of the purpose of the go file.

RESTRICTIONS:
    AWS functions:
    * Program must have access to a .aws/credentials file in the default location.
    * This will only access system parameters that start with '/sote' (ROOTPATH).
    * {Enter other restrictions here for AWS

    {Other catagories of restrictions}
    * {List of restrictions for the catagory

NOTES:
    {Enter any additional notes that you believe will help the next developer.}

COPYRIGHT:
	Copyright 2022
	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.

*/
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func signalHandler(signal os.Signal) {
	fmt.Printf("\nCaught signal: %+v", signal)
	fmt.Println("\nWait for 1 second to finish processing")
	time.Sleep(1 * time.Second)

	switch signal {
	case syscall.SIGHUP: // kill -SIGHUP XXXX
		fmt.Println("- got hungup")

	case syscall.SIGINT: // kill -SIGINT XXXX or Ctrl+c
		fmt.Println("- got ctrl+c")

	case syscall.SIGTERM: // kill -SIGTERM XXXX
		fmt.Println("- got force stop")

	case syscall.SIGQUIT: // kill -SIGQUIT XXXX
		fmt.Println("- stop and core dump")

	default:
		fmt.Println("- unknown signal")
	}

	fmt.Println("\nFinished server cleanup")
	os.Exit(0)
}

func initSignals() {
	var captureSignal = make(chan os.Signal, 1)
	signal.Notify(captureSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	signalHandler(<-captureSignal)
}

func main() {

	var (
		totalSleep = 0
	)

	go initSignals()

	for true {
		time.Sleep(1 * time.Minute)
		totalSleep++
		log.Printf("Total Sleep time: %v\n", totalSleep)
	}
}
