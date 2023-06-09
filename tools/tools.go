//go:build tools
// +build tools

package tools

import (
	// Blank import ensures that these repos are included in go.mod, so we can
	// build CLI tools from the ./vendor/ directory during `make install_deps`
	_ "github.com/Clever/launch-gen"
	_ "github.com/cespare/reflex"
	_ "github.com/get-woke/woke"
	_ "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/kyleconroy/sqlc/cmd/sqlc"
	_ "github.com/nomad-software/vend"
)
