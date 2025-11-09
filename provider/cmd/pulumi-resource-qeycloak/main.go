// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"
	"strings"

	p "github.com/pulumi/pulumi-go-provider"

	aci "github.com/netascode/pulumi-aci/provider/pkg/provider"
	"github.com/netascode/pulumi-aci/provider/pkg/version"
)

// A provider is a program that listens for requests from the Pulumi engine
// to interact with cloud providers using a CRUD-based model.
func main() {
	version := version.Version
	if strings.HasPrefix(version, "v") {
		version = version[1:]
	}

	// This method defines the provider implemented in this repository.
	aciProvider := aci.NewProvider()

	// This method starts serving requests using the ACI provider.
	err := p.RunProvider("aci", version, aciProvider)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}
