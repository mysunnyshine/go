// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package load

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
)

// DebugDeprecatedImportcfg is installed as the undocumented -debug-deprecated-importcfg build flag.
// It is useful for debugging subtle problems in the go command logic but not something
// we want users to depend on. The hope is that the "deprecated" will make that clear.
// We intend to remove this flag in Go 1.11.
var DebugDeprecatedImportcfg debugDeprecatedImportcfgFlag

type debugDeprecatedImportcfgFlag struct {
	enabled bool
	pkgs    map[string]*debugDeprecatedImportcfgPkg
}

type debugDeprecatedImportcfgPkg struct {
	Dir    string
	Import map[string]string
}

var (
	debugDeprecatedImportcfgMagic = []byte("# debug-deprecated-importcfg\n")
	errImportcfgSyntax            = errors.New("malformed syntax")
)

func (f *debugDeprecatedImportcfgFlag) String() string { return "" }

func (f *debugDeprecatedImportcfgFlag) Set(x string) error {
	if x == "" {
		*f = debugDeprecatedImportcfgFlag{}
		return nil
	}
	data, err := ioutil.ReadFile(x)
	if err != nil {
		return err
	}

	if !bytes.HasPrefix(data, debugDeprecatedImportcfgMagic) {
		return errImportcfgSyntax
	}
	data = data[len(debugDeprecatedImportcfgMagic):]

	f.pkgs = nil
	if err := json.Unmarshal(data, &f.pkgs); err != nil {
		return errImportcfgSyntax
	}
	f.enabled = true
	return nil
}

func (f *debugDeprecatedImportcfgFlag) lookup(parent *Package, path string) (dir, newPath string) {
	if parent == nil {
		if p1 := f.pkgs[path]; p1 != nil {
			return p1.Dir, path
		}
		return "", ""
	}
	if p1 := f.pkgs[parent.ImportPath]; p1 != nil {
		if newPath := p1.Import[path]; newPath != "" {
			if p2 := f.pkgs[newPath]; p2 != nil {
				return p2.Dir, newPath
			}
		}
	}
	return "", ""
}
