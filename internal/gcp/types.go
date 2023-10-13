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
