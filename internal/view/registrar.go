package view

func loadCustomViewers() MetaViewers {
	m := make(MetaViewers, 5)
	coreViewers(m)
	return m
}

func coreViewers(vv MetaViewers) {
	// TODO create consts instead of hardcoded
	vv["ec2"] = MetaViewer{
		viewerFn: NewEC2,
	}
	vv["s3"] = MetaViewer{
		viewerFn: NewS3,
	}
	vv["sg"] = MetaViewer{
		viewerFn: NewSG,
	}
	vv["iam:u"] = MetaViewer{
		viewerFn: NewIAMU,
	}
	vv["ebs"] = MetaViewer{
		viewerFn: NewEBS,
	}
	vv["iam:g"] = MetaViewer{
		viewerFn: NewIAMUG,
	}
}
