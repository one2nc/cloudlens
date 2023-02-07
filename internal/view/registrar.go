package view

func loadCustomViewers() MetaViewers {
	m := make(MetaViewers, 5)
	coreViewers(m)
	return m
}

func coreViewers(vv MetaViewers) {
	vv["ec2"] = MetaViewer{
		viewerFn: NewEC2,
	}
	vv["s3"] = MetaViewer{
		viewerFn: NewS3,
	}
	vv["sg"] = MetaViewer{
		viewerFn: NewSG,
	}
}