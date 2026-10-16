package md

import _ "embed"

var (
	//go:embed Provider.md
	Provider string
	//go:embed ProviderUsername.md
	ProviderUsername string
	//go:embed ProviderPassword.md
	ProviderPassword string
	//go:embed ResourceRecordA.md
	ResourceRecordA string
	//go:embed ResourceRecordZone.md
	ResourceRecordZone string
	//go:embed SoaMinimumTtl.md
	SoaMinimumTtl string
	//go:embed SoaTtl.md
	SoaTtl string
)
