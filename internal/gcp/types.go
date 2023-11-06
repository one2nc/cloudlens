package gcp

type StorageResp struct {
	BucketName   string
	CreationTime string
	Region       string
}

type StorageObjResp struct {
	SizeInBytes                                        int64
	Name, ObjectType, LastModified, Size, StorageClass string
}

type VMResp struct {
	// Instance         .Instance
	InstanceId       string
	InstanceType     string
	AvailabilityZone string
	InstanceState    string
	PublicDNS        string
	MonitoringState  string
	LaunchTime       string
}

type DiskResp struct {
	Name, Type, Size,CreationTime,Status string
	Zone string
}
