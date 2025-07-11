/*
Copyright © 2023 IceWhaleTech

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
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/BeesNestInc/CassetteOS-CLI/codegen/app_management"
	"github.com/spf13/cobra"
)

// appManagementUnregisterAppStoreCmd represents the appManagementUnregisterAppStore command
var appManagementUnregisterAppStoreCmd = &cobra.Command{
	Use:     "app-store <id>",
	Short:   "unregister an app store by id",
	Aliases: []string{"appstore"},
	Args: cobra.MatchAll(cobra.ExactArgs(1), func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil || id < 0 {
			return fmt.Errorf("id must be a number larger or equal to 0")
		}

		return nil
	}),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootURL, err := rootCmd.PersistentFlags().GetString(FlagRootURL)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("http://%s/%s", rootURL, BasePathAppManagement)

		appStoreID, err := strconv.Atoi(cmd.Flags().Arg(0))
		if err != nil || appStoreID < 0 {
			return fmt.Errorf("how can it get here?? should have been validated in cobra.MatchAll(...)")
		}

		client, err := app_management.NewClientWithResponses(url)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()

		response, err := client.UnregisterAppStoreWithResponse(ctx, appStoreID)
		if err != nil {
			return err
		}

		var baseResponse app_management.BaseResponse

		if response.StatusCode() != http.StatusOK {
			if err := json.Unmarshal(response.Body, &baseResponse); err != nil {
				return fmt.Errorf("%s - %s", response.Status(), response.Body)
			}

			return fmt.Errorf("%s - %s", response.Status(), *baseResponse.Message)
		}

		if err := json.Unmarshal(response.Body, &baseResponse); err != nil {
			return fmt.Errorf("%s - %s", response.Status(), response.Body)
		}

		fmt.Println(*baseResponse.Message)

		return nil
	},
}

func init() {
	appManagementUnregisterCmd.AddCommand(appManagementUnregisterAppStoreCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appManagementUnregisterAppStoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appManagementUnregisterAppStoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
