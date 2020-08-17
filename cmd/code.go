/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

// signinCmd represents the signin command
var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Authorize the application from your account",
	Long:  `Helps to generate the code by providing the access, we promise not to pose without your consent`,
	Run: func(cmd *cobra.Command, args []string) {

		path, _ := os.LookupEnv("HOME")
		f, err := os.Create(path + "/.producthunt")

		if err != nil {
			panic(err)
		}
		defer f.Close()

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
		s.Start()
		url := "https://api.producthunt.com/v2/oauth/token"
		values := map[string]string{"client_id": "K0G_mZKnQkvmBPTHXC7bKAXgJZlLzgA0TePqFpn2yJU", "client_secret": "4Upegb9eVYv6PzdInfXEzXP5jSX96KWyQ61-zMxL6Ug", "grant_type": "authorization_code", "redirect_uri": "https://producthuntcli.netlify.app/", "code": args[0]}
		jsonValue, _ := json.Marshal(values)
		request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
		request.Header.Add("Accept", "application/json")
		request.Header.Add("Host", "api.producthunt.com")
		request.Header.Add("Content-Type", "application/json")
		client := &http.Client{Timeout: time.Second * 30}
		response, err := client.Do(request)
		s.Stop()
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}
		data, _ := ioutil.ReadAll(response.Body)
		var result map[string]interface{}

		// Unmarshal or Decode the JSON to the interface.
		json.Unmarshal([]byte(data), &result)
		if result["access_token"] == nil {
			fmt.Println("Please signin again to generate the fresh code by running, `producthunt signin`")
		}

		_, _ = f.WriteString(fmt.Sprintf("%s", result["access_token"]))
		fmt.Println("Code added successfully!")
	},
}

func init() {

	codeCmd.Flags().StringP("auth", "a", "", "Auth code (required)")
	rootCmd.AddCommand(codeCmd)
}
