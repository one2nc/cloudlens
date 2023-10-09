package view

import "github.com/one2nc/cloudlens/internal"

func loadCustomViewers() MetaViewers {
	m := make(MetaViewers, 5)
	coreViewers(m)
	return m
}

func coreViewers(vv MetaViewers) {
	// TODO create consts instead of hardcoded
	vv[internal.LowercaseEc2] = MetaViewer{
		viewerFn: NewEC2,
	}
	vv[internal.LowercaseS3] = MetaViewer{
		viewerFn: NewS3,
	}
	vv[internal.LowercaseSg] = MetaViewer{
		viewerFn: NewSG,
	}
	vv[internal.UppercaseSg] = MetaViewer{
		viewerFn: NewSG,
	}
	vv[internal.LowercaseIamUser] = MetaViewer{
		viewerFn: NewIAMU,
	}
	vv[internal.LowercaseEBS] = MetaViewer{
		viewerFn: NewEBS,
	}
	vv[internal.LowercaseIamGroup] = MetaViewer{
		viewerFn: NewIAMUG,
	}
	vv[internal.LowercaseIamRole] = MetaViewer{
		viewerFn: NewIamRole,
	}
	vv[internal.LowercaseEc2Snapshot] = MetaViewer{
		viewerFn: NewEC2S,
	}
	vv[internal.LowercaseEc2Image] = MetaViewer{
		viewerFn: NewEC2I,
	}
	vv[internal.LowercaseSQS] = MetaViewer{
		viewerFn: NewSQS,
	}
	vv[internal.LowercaseVPC] = MetaViewer{
		viewerFn: NewVPC,
	}
	vv[internal.LowercaseSubnet] = MetaViewer{
		viewerFn: NewSubnet,
	}
	vv[internal.LowercaseLamda] = MetaViewer{
		viewerFn: NewLambda,
	}
	vv[internal.LowercaseStorage] = MetaViewer{
		viewerFn: NewStorage,
	}
}
