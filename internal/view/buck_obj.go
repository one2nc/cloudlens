package view

import (
	"context"
	"fmt"

	"github.com/atotto/clipboard"
	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/one2nc/cloudlens/internal/ui"
	"github.com/rs/zerolog/log"
)

type S3FileViewer struct {
	name, path string
	ResourceViewer
}

func NewS3FileViewer(path, resource string) *S3FileViewer {
	var obj S3FileViewer
	obj.name = resource
	obj.path = path
	if obj.path == "s3://" {
		obj.path = obj.path + obj.name
	}
	obj.ResourceViewer = NewBrowser("OBJ")
	obj.AddBindKeysFn(obj.bindKeys)
	//s3.GetTable().SetEnterFn(s3.describeInstace)
	return &obj
}

func (obj *S3FileViewer) Name() string {
	return obj.name
}

func (obj *S3FileViewer) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftN:    ui.NewKeyAction("Sort Name", obj.GetTable().SortColCmd("Name", true), true),
		ui.KeyShiftM:    ui.NewKeyAction("Sort Modification-Time", obj.GetTable().SortColCmd("Last-Modified", true), true),
		ui.KeyShiftS:    ui.NewKeyAction("Sort Size", obj.GetTable().SortColCmd("Size", true), true),
		ui.KeyShiftC:    ui.NewKeyAction("Sort Storage-Class", obj.GetTable().SortColCmd("Storage-Class", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", obj.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", obj.enterCmd, false),
		tcell.KeyCtrlD:  ui.NewKeyAction("Download Object", obj.downloadCmd, true),
		tcell.KeyCtrlP:  ui.NewKeyAction("Pre-Signed URL", obj.preSignedUrlCmd, true),
	})
}

func (obj *S3FileViewer) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	objName := obj.GetTable().GetSelectedItem()
	fileType := obj.GetTable().GetSecondColumn()
	if fileType == "Folder" {
		o := NewS3FileViewer(obj.path+"/"+objName, objName)
		ctx := obj.App().GetContext()
		bn := ctx.Value(internal.BucketName)
		fn := fmt.Sprintf("%v%v/", ctx.Value(internal.FolderName), objName)
		log.Info().Msg(fmt.Sprintf("In view Folder Name: %v", fn))
		ctx = context.WithValue(obj.App().context, internal.BucketName, bn)
		obj.App().SetContext(ctx)
		ctx = context.WithValue(obj.App().context, internal.FolderName, fn)
		obj.App().SetContext(ctx)

		obj.App().Flash().Info(fmt.Sprintf("Bucket Name: %v", bn))
		obj.App().inject(o)
		o.GetTable().SetTitle(o.path)
	}

	return evt
}

func (obj *S3FileViewer) downloadCmd(evt *tcell.EventKey) *tcell.EventKey {
	objName := obj.GetTable().GetSelectedItem()
	fileType := obj.GetTable().GetSecondColumn()

	if fileType == "File" {
		ctx := obj.App().GetContext()
		op := getObjectParams(ctx, objName)
		res := aws.DownloadObject(op.cfg, op.bucketName, op.key)
		obj.App().Flash().Info(res)
	}

	return nil
}

func (obj *S3FileViewer) preSignedUrlCmd(evt *tcell.EventKey) *tcell.EventKey {
	objNmae := obj.GetTable().GetSelectedItem()
	fileType := obj.GetTable().GetSecondColumn()

	if fileType == "File" {
		ctx := obj.App().GetContext()
		op := getObjectParams(ctx, objNmae)
		url := aws.GetPreSignedUrl(op.cfg, op.bucketName, op.key)
		log.Info().Msg(fmt.Sprintf("In view Presigned URL: %v", url))
		clipboard.WriteAll(url)
		obj.App().Flash().Info("Presigned URL Copied to Clipboard.")
	}

	return nil
}

func getObjectParams(ctx context.Context, objName string) ObjectParams {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}

	bn := fmt.Sprintf("%v", ctx.Value(internal.BucketName))
	fn := fmt.Sprintf("%v", ctx.Value(internal.FolderName))
	log.Info().Msg(fmt.Sprintf("In view Bucket Name: %v", bn))
	log.Info().Msg(fmt.Sprintf("In view Folder Name: %v", fn))
	log.Info().Msg(fmt.Sprintf("In view Object Name: %v", objName))
	key := fn + objName
	log.Info().Msg(fmt.Sprintf("In view key: %v", key))
	return ObjectParams{
		cfg:        cfg,
		bucketName: bn,
		key:        key,
	}
}
