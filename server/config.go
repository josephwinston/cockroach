// Copyright 2014 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Bram Gruneir (bram.gruneir@gmail.com)

package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/cockroachdb/cockroach/util"
	"github.com/cockroachdb/cockroach/util/log"
)

// sendAdminRequest send an HTTP request and processes the response for
// its body or error message if a non-200 response code.
func sendAdminRequest(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, util.Errorf("admin REST request failed: %s", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, util.Errorf("unable to read admin REST response: %s", err)
	}
	if resp.StatusCode != 200 {
		return nil, util.Errorf("%s: %s", resp.Status, string(b))
	}
	return b, nil
}

// Gets a friendly name for output based on the passed in config prefix.
func getFriendlyNameFromPrefix(prefix string) string {
	switch prefix {
	case acctPathPrefix:
		return "accounting"
	case permPathPrefix:
		return "permission"
	case zonePathPrefix:
		return "zone"
	default:
		return "unknown"
	}
}

// runGetConfig invokes the REST API with GET action and key prefix as path.
func runGetConfig(ctx *Context, prefix, keyPrefix string) {
	friendlyName := getFriendlyNameFromPrefix(prefix)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s://%s%s/%s", adminScheme, ctx.Addr, prefix, keyPrefix), nil)
	if err != nil {
		log.Errorf("unable to create request to admin REST endpoint: %s", err)
		return
	}
	req.Header.Add("Accept", "text/yaml")
	// TODO(spencer): need to move to SSL.
	b, err := sendAdminRequest(req)
	if err != nil {
		log.Errorf("admin REST request failed: %s", err)
		return
	}
	fmt.Fprintf(os.Stdout, "%s config for key prefix %q:\n%s\n", friendlyName, keyPrefix, string(b))
}

// RunGetAcct gets the account from the given key.
func RunGetAcct(ctx *Context, keyPrefix string) {
	runGetConfig(ctx, acctPathPrefix, keyPrefix)
}

// RunGetPerm gets the permission from the given key.
func RunGetPerm(ctx *Context, keyPrefix string) {
	runGetConfig(ctx, permPathPrefix, keyPrefix)
}

// RunGetZone gets the zone from the given key.
func RunGetZone(ctx *Context, keyPrefix string) {
	runGetConfig(ctx, zonePathPrefix, keyPrefix)
}

// runLsConfigs invokes the REST API with GET action and no path, which
// fetches a list of all configuration prefixes.
// The type of config that is listed is based on the passed in prefix.
// The optional regexp is applied to the complete list and matching prefixes
// displayed.
func runLsConfigs(ctx *Context, prefix, pattern string) {
	friendlyName := getFriendlyNameFromPrefix(prefix)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s://%s%s", adminScheme, ctx.Addr, prefix), nil)
	if err != nil {
		log.Errorf("unable to create request to admin REST endpoint: %s", err)
		return
	}
	b, err := sendAdminRequest(req)
	if err != nil {
		log.Errorf("admin REST request failed: %s", err)
		return
	}
	var prefixes []string
	if err = json.Unmarshal(b, &prefixes); err != nil {
		log.Errorf("unable to parse admin REST response: %s", err)
		return
	}
	var re *regexp.Regexp
	if len(pattern) > 0 {
		if re, err = regexp.Compile(pattern); err != nil {
			log.Warningf("invalid regular expression %q; skipping regexp match and listing all %s prefixes", pattern, friendlyName)
			re = nil
		}
	}
	for _, prefix := range prefixes {
		if re != nil {
			unescaped, err := url.QueryUnescape(prefix)
			if err != nil || !re.MatchString(unescaped) {
				continue
			}
		}
		if prefix == "" {
			prefix = "[default]"
		}
		fmt.Fprintf(os.Stdout, "%s\n", prefix)
	}
}

// RunLsAcct lists accounts.
func RunLsAcct(ctx *Context, pattern string) {
	runLsConfigs(ctx, acctPathPrefix, pattern)
}

// RunLsPerm lists permissions.
func RunLsPerm(ctx *Context, pattern string) {
	runLsConfigs(ctx, permPathPrefix, pattern)
}

// RunLsZone lists zones.
func RunLsZone(ctx *Context, pattern string) {
	runLsConfigs(ctx, zonePathPrefix, pattern)
}

// runRmConfig invokes the REST API with DELETE action and key prefix as path.
// The type of config that is removed is based on the passed in prefix.
func runRmConfig(ctx *Context, prefix, keyPrefix string) {
	friendlyName := getFriendlyNameFromPrefix(prefix)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s://%s%s/%s", adminScheme, ctx.Addr, prefix, keyPrefix), nil)
	if err != nil {
		log.Errorf("unable to create request to admin REST endpoint: %s", err)
		return
	}
	// TODO(spencer): need to move to SSL.
	_, err = sendAdminRequest(req)
	if err != nil {
		log.Errorf("admin REST request failed: %s", err)
		return
	}
	fmt.Fprintf(os.Stdout, "removed %s config for key prefix %q\n", friendlyName, keyPrefix)
}

// RunRmAcct removes the account with the given key.
func RunRmAcct(ctx *Context, keyPrefix string) {
	runRmConfig(ctx, acctPathPrefix, keyPrefix)
}

// RunRmPerm removes the permission with the given key.
func RunRmPerm(ctx *Context, keyPrefix string) {
	runRmConfig(ctx, permPathPrefix, keyPrefix)
}

// RunRmZone removes the zone with the given key.
func RunRmZone(ctx *Context, keyPrefix string) {
	runRmConfig(ctx, zonePathPrefix, keyPrefix)
}

// runSetConfig invokes the REST API with POST action and key prefix as
// path. The specified configuration file is read from disk and sent
// as the POST body.
// The type of config that is set is based on the passed in prefix.
func runSetConfig(ctx *Context, prefix, keyPrefix, configFileName string) {
	friendlyName := getFriendlyNameFromPrefix(prefix)
	// Read in the config file.
	body, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Errorf("unable to read %s config file %q: %s", friendlyName, configFileName, err)
		return
	}
	// Send to admin REST API.
	req, err := http.NewRequest("POST", fmt.Sprintf("%s://%s%s/%s", adminScheme, ctx.Addr, prefix, keyPrefix), bytes.NewReader(body))
	if err != nil {
		log.Errorf("unable to create request to admin REST endpoint: %s", err)
		return
	}
	req.Header.Add("Content-Type", "text/yaml")
	// TODO(spencer): need to move to SSL.
	_, err = sendAdminRequest(req)
	if err != nil {
		log.Errorf("admin REST request failed: %s", err)
		return
	}
	fmt.Fprintf(os.Stdout, "set %s config for key prefix %q\n", friendlyName, keyPrefix)
}

// RunSetAcct sets the account to the key given the yaml filename.
func RunSetAcct(ctx *Context, keyPrefix, configFileName string) {
	runSetConfig(ctx, acctPathPrefix, keyPrefix, configFileName)
}

// RunSetPerm sets the permission to the key given the yaml filename.
func RunSetPerm(ctx *Context, keyPrefix, configFileName string) {
	runSetConfig(ctx, permPathPrefix, keyPrefix, configFileName)
}

// RunSetZone sets the zone to the key given the yaml filename.
func RunSetZone(ctx *Context, keyPrefix, configFileName string) {
	runSetConfig(ctx, zonePathPrefix, keyPrefix, configFileName)
}
