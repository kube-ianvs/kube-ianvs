/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

func NewGatewayCmd() (gatewayCmd *cobra.Command) {

	var addr *string
	var certFile *string
	var keyFile *string
	gatewayCmd = &cobra.Command{
		Use:   "gateway",
		Short: "ianvs apiserver gateway",
		Long:  "ianvs apiserver gateway",
		Run: func(cmd *cobra.Command, args []string) {
			http.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
				io.WriteString(w, "hello, world!\n")
			})
			fmt.Println(certFile, keyFile)
			if e := http.ListenAndServeTLS(*addr, *certFile, *keyFile, nil); e != nil {
				log.Fatal("ListenAndServe: ", e)
			}
		},
	}

	addr = gatewayCmd.PersistentFlags().String("addr", "", "https listen address")
	certFile = gatewayCmd.PersistentFlags().String("certFile", "", "https cert file")
	keyFile = gatewayCmd.PersistentFlags().String("keyFile", "", "https key file")
	gatewayCmd.MarkPersistentFlagRequired("addr")
	gatewayCmd.MarkPersistentFlagRequired("certFile")
	gatewayCmd.MarkPersistentFlagRequired("keyFile")
	return
}
