// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wptdashboard

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestParseSHAParam(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/", nil)
	runSHA, err := ParseSHAParam(r)
	assert.Nil(t, err)
	assert.Equal(t, "latest", runSHA)
}

func TestParseSHAParam_2(t *testing.T) {
	sha := "0123456789"
	r := httptest.NewRequest("GET", "http://wpt.fyi/?sha="+sha, nil)
	runSHA, err := ParseSHAParam(r)
	assert.Nil(t, err)
	assert.Equal(t, sha, runSHA)
}

func TestParseSHAParam_BadRequest(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/?sha=%zz", nil)
	runSHA, err := ParseSHAParam(r)
	assert.NotNil(t, err)
	assert.Equal(t, "latest", runSHA)
}

func TestParseSHAParam_NonSHA(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/?sha=123", nil)
	runSHA, err := ParseSHAParam(r)
	assert.Nil(t, err)
	assert.Equal(t, "latest", runSHA)
}

func TestParseSHAParam_NonSHA_2(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/?sha=zapper0123", nil)
	runSHA, err := ParseSHAParam(r)
	assert.Nil(t, err)
	assert.Equal(t, "latest", runSHA)
}

func TestParseBrowserParam(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/", nil)
	browser, err := ParseBrowserParam(r)
	assert.Nil(t, err)
	assert.Equal(t, "", browser)
}

func TestParseBrowserParam_Chrome(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/?browser=chrome", nil)
	browser, err := ParseBrowserParam(r)
	assert.Nil(t, err)
	assert.Equal(t, "chrome", browser)
}

func TestParseBrowserParam_Invalid(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/?browser=invalid", nil)
	browser, err := ParseBrowserParam(r)
	assert.NotNil(t, err)
	assert.Equal(t, "", browser)
}

func TestParseBrowsersParam(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/", nil)
	browsers, err := ParseBrowsersParam(r)
	assert.Nil(t, err)
	defaultBrowsers, err := GetBrowserNames()
	assert.Equal(t, defaultBrowsers, browsers)
}

func TestParseBrowsersParam_ChromeSafari(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/?browsers=chrome,safari", nil)
	browsers, err := ParseBrowsersParam(r)
	assert.Nil(t, err)
	assert.Equal(t, []string{"chrome", "safari"}, browsers)
}

func TestParseBrowsersParam_ChromeInvalid(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/?browsers=chrome,invalid", nil)
	browsers, err := ParseBrowsersParam(r)
	assert.Nil(t, err)
	assert.Equal(t, []string{"chrome"}, browsers)
}

func TestParseBrowsersParam_AllInvalid(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/?browsers=notabrowser,invalid", nil)
	browsers, err := ParseBrowsersParam(r)
	assert.Nil(t, err)
	assert.Empty(t, browsers)
}

func TestParseBrowsersParam_EmptyCommas(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/?browsers=,notabrowser,,,,invalid,,", nil)
	browsers, err := ParseBrowsersParam(r)
	assert.Nil(t, err)
	assert.Empty(t, browsers)
}

func TestParseBrowsersParam_SafariChrome(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/?browsers=safari,chrome", nil)
	browsers, err := ParseBrowsersParam(r)
	assert.Nil(t, err)
	assert.Equal(t, []string{"chrome", "safari"}, browsers)
}

func TestParseBrowsersParam_MultiBrowserParam_SafariChrome(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/?browser=safari&browser=chrome", nil)
	browsers, err := ParseBrowsersParam(r)
	assert.Nil(t, err)
	assert.Equal(t, []string{"chrome", "safari"}, browsers)
}

func TestParseBrowsersParam_MultiBrowserParam_SafariInvalid(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/?browser=safari&browser=invalid", nil)
	browsers, err := ParseBrowsersParam(r)
	assert.Nil(t, err)
	assert.Equal(t, []string{"safari"}, browsers)
}

func TestParseBrowsersParam_MultiBrowserParam_AllInvalid(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/?browser=invalid&browser=notabrowser", nil)
	browsers, err := ParseBrowsersParam(r)
	assert.Nil(t, err)
	assert.Empty(t, browsers)
}

func TestParsePathsParam_Missing(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/api/diff", nil)
	paths := ParsePathsParam(r)
	assert.Nil(t, paths)
}

func TestParsePathsParam_Path_Duplicate(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/api/diff?path=/css&path=/css", nil)
	paths := ParsePathsParam(r)
	assert.Len(t, paths.ToSlice(), 1)
}

func TestParsePathsParam_Paths_Duplicate(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/api/diff?paths=/css,/css", nil)
	paths := ParsePathsParam(r)
	assert.Len(t, paths.ToSlice(), 1)
}

func TestParsePathsParam_PathsAndPath_Duplicate(t *testing.T) {
	r := httptest.NewRequest("GET", "http://wpt.fyi/api/diff?paths=/css&path=/css", nil)
	paths := ParsePathsParam(r)
	assert.Len(t, paths.ToSlice(), 1)
}