package view

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal"
	"github.com/one2nc/cloud-lens/internal/ui"
	"github.com/rs/zerolog/log"
)

type BObj struct {
	ResourceViewer
}

func NewS3FileViewer(resource string) ResourceViewer {
	var obj BObj
	obj.ResourceViewer = NewBrowser(resource)
	obj.AddBindKeysFn(obj.bindKeys)
	//s3.GetTable().SetEnterFn(s3.describeInstace)
	return &obj
}
func (obj *BObj) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftN:    ui.NewKeyAction("Sort Name", obj.GetTable().SortColCmd("Name", true), true),
		ui.KeyShiftM:    ui.NewKeyAction("Sort Modification-Time", obj.GetTable().SortColCmd("Last-Modified", true), true),
		ui.KeyShiftS:    ui.NewKeyAction("Sort Size", obj.GetTable().SortColCmd("Size", true), true),
		ui.KeyShiftC:    ui.NewKeyAction("Sort Storage-Class", obj.GetTable().SortColCmd("Storage-Class", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", obj.App().PrevCmd, true),
		tcell.KeyEnter:  ui.NewKeyAction("View", obj.enterCmd, true),
	})
}

func (obj *BObj) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	oName := obj.GetTable().GetSelectedItem()
	fileType := obj.GetTable().GetSecondColumn()
	if fileType == "Folder" {
		o := NewS3FileViewer("OBJ")
		ctx := obj.App().GetContext()
		bn := ctx.Value(internal.BucketName)
		fn := fmt.Sprintf("%v%v/", ctx.Value(internal.FolderName), oName)
		log.Info().Msg(fmt.Sprintf("In view Folder Name: %v", fn))
		ctx = context.WithValue(obj.App().context, internal.BucketName, bn)
		obj.App().SetContext(ctx)
		ctx = context.WithValue(obj.App().context, internal.FolderName, fn)
		obj.App().SetContext(ctx)

		obj.App().Flash().Info("Bucket Name: " + oName)
		// println(bName)
		obj.App().inject(o)
	}
	return nil
}
