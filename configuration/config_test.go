// Copyright (c) 2016 The GOST Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExternalUri(t *testing.T) {
	cfg := Config{}
	cfg.Server.ExternalURI = "http://test.com/"
	assert.Equal(t, "http://test.com", cfg.GetExternalServerURI(), "Trailing slash not removed by GetExternalServerUri")
}

func TestGetInternalUri(t *testing.T) {
	cfg := Config{}
	cfg.Server.Host = "localhost"
	cfg.Server.Port = 8080
	assert.Equal(t, "localhost:8080", cfg.GetInternalServerURI(), "Internal server uri not constructed correctly based on config server host and port")
}
