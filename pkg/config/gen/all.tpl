// Copyright (C) 2026 XLR8discovery PBC
// See LICENSE for copying information.

package config

func All() map[string][]Option {
	return map[string][]Option{ {{ range $name, $def := .Configs }}
		"{{ $name }}": {{ $name | goName }}Config(),{{ end }}
	}
}
