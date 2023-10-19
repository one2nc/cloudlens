package view

import (
	"context"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/gcp"
	"github.com/one2nc/cloudlens/internal/ui"
	"github.com/rs/zerolog/log"
)

type StorageFileViewer struct {
	// name, path string
	folderName string
	path       string
	bucketName string
	ResourceViewer
}

func NewStorageFileViewer(path string, bucketName string, folderName string) *StorageFileViewer {
	var obj StorageFileViewer

	obj.path = path
	obj.bucketName = bucketName
	obj.folderName = folderName
	obj.ResourceViewer = NewBrowser(internal.StorageObject)
	obj.AddBindKeysFn(obj.bindKeys)
	return &obj
}

func (obj *StorageFileViewer) Name() string {

	if obj.folderName != "" {
		return strings.ReplaceAll(obj.folderName, "/", "")
	}
	return obj.bucketName
}

func (obj *StorageFileViewer) bindKeys(aa ui.KeyActions) {
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

func (obj *StorageFileViewer) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	objName := obj.GetTable().GetSelectedItem()
	fileType := obj.GetTable().GetSecondColumn()
	if fileType == internal.FOLDER_TYPE {
		ctx := obj.App().GetContext()
		bn := ctx.Value(internal.BucketName).(string)
		o := NewStorageFileViewer(obj.path+objName, bn, objName)

		log.Info().Msg(fmt.Sprintf("In view Folder Name: %v", objName))
		ctx = context.WithValue(obj.App().context, internal.BucketName, bn)
		obj.App().SetContext(ctx)
		ctx = context.WithValue(obj.App().context, internal.FolderName, o.path)
		obj.App().SetContext(ctx)

		obj.App().Flash().Info(fmt.Sprintf("Bucket Name: %v", bn))
		obj.App().inject(o)
		o.GetTable().SetTitle(o.path)
	}

	return evt
}

func (obj *StorageFileViewer) downloadCmd(evt *tcell.EventKey) *tcell.EventKey {
	objName := obj.GetTable().GetSelectedItem()
	fileType := obj.GetTable().GetSecondColumn()
	ctx := obj.App().GetContext()

	if fileType == internal.FILE_TYPE {
		res := gcp.DownloadObject(ctx, obj.bucketName, obj.path, objName)
		obj.App().Flash().Info(res)
	}

	return nil
}

func (obj *StorageFileViewer) preSignedUrlCmd(evt *tcell.EventKey) *tcell.EventKey {
	objName := obj.GetTable().GetSelectedItem()
	fileType := obj.GetTable().GetSecondColumn()
	ctx := obj.App().GetContext()

	if fileType == internal.FILE_TYPE {
		url := gcp.GetPreSignedUrl(ctx, obj.bucketName, obj.path, objName)
		if url != "" {
			log.Info().Msg(fmt.Sprintf("In view Presigned URL: %v", url))
			clipboard.WriteAll(url)
			obj.App().Flash().Info("Presigned URL Copied to Clipboard.")
		}
	}

	return nil
}
