package md

import _ "embed"

var (
	//go:embed Provider.md
	Provider string
	//go:embed ProviderUsername.md
	ProviderUsername string
	//go:embed ProviderPassword.md
	ProviderPassword string
	//go:embed RecordPriority.md
	RecordPriority string
	//go:embed RecordZone.md
	RecordZone string
	//go:embed ResourceRecordA.md
	ResourceRecordA string
	//go:embed ResourceRecordAAAA.md
	ResourceRecordAAAA string
	//go:embed ResourceRecordCNAME.md
	ResourceRecordCNAME string
	//go:embed ResourceRecordCNAMECanonical.md
	ResourceRecordCNAMECanonical string
	//go:embed ResourceRecordMX.md
	ResourceRecordMX string
	//go:embed ResourceRecordMXMailServer.md
	ResourceRecordMXMailServer string
	//go:embed ResourceRecordTXT.md
	ResourceRecordTXT string
	//go:embed ResourceRecordTXTText.md
	ResourceRecordTXTText string
	//go:embed SoaMinimumTtl.md
	SoaMinimumTtl string
	//go:embed SoaTtl.md
	SoaTtl string
)
