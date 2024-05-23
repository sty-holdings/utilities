// Package main.go
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

// import (
// 	"log"
//
// 	"albert/constants"
// 	"albert/core"
// 	"albert/core/coreHelpers"
// 	"firebase.google.com/go/auth"
// )
//
// const (
// // Add Constants to the constants.go file
// )
//
// // Add types to the types.go file
//
// var (
// // Add Variables here for the file (Remember, they are global)
// )
//
// func init() {
// 	// Set up goes here
// }
//
// func main() {
//
// 	var (
// 		err       error
// 		tFirebase coreHelpers.FirebaseFirestoreHelper
// 		tUser     *auth.UserRecord
// 	)
//
// 	tFirebase.AppPtr, tFirebase.AuthPtr, tFirebase.FirestoreClientPtr, err = core.GetFirebaseFirestoreConnection(constants.TEST_FIREBASE_CREDENTIALS)
//
// 	params := (&auth.UserToCreate{}).
// 		Email("user@example.com").
// 		EmailVerified(false).
// 		PhoneNumber("+16504830648").
// 		Password("secretPassword").
// 		DisplayName("Donald Drump").
// 		PhotoURL("http://www.example.com/12345678/photo.png").
// 		Disabled(false)
// 	tUser, err = tFirebase.AuthPtr.CreateUser(core.CTXBackground, params)
// 	if err != nil {
// 		log.Fatalf("error creating user: %v\n", err)
// 	}
// 	log.Printf("Successfully created user: %v\n", tUser)
//
// 	tFirebase.AuthPtr.DeleteUser(core.CTXBackground, tUser.UID)
// }
