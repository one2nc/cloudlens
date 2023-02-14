package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloud-lens/internal/ui"
)

type SQS struct {
	ResourceViewer
}

// NewPod returns a new viewer.
func NewSQS(resource string) ResourceViewer {
	var sqs SQS
	sqs.ResourceViewer = NewBrowser(resource)
	sqs.AddBindKeysFn(sqs.bindKeys)
	// e.GetTable().SetEnterFn(e.describeInstace)
	return &sqs
}

func (sqs *SQS) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftN:    ui.NewKeyAction("Sort Name", sqs.GetTable().SortColCmd("Name", true), true),
		ui.KeyShiftT:    ui.NewKeyAction("Sort Type", sqs.GetTable().SortColCmd("Type", true), true),
		ui.KeyShiftC:    ui.NewKeyAction("Sort Created", sqs.GetTable().SortColCmd("Created", true), true),
		ui.KeyShiftM:    ui.NewKeyAction("Sort Messages-Available", sqs.GetTable().SortColCmd("Messages-Available", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", sqs.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", sqs.enterCmd, false),
	})
}
func (sqs *SQS) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	queueUrl := sqs.GetTable().GetSelectedItem()
	if queueUrl != "" {
		f := describeResource
		if sqs.GetTable().enterFn != nil {
			f = sqs.GetTable().enterFn
		}
		f(sqs.App(), sqs.GetTable().GetModel(), sqs.Resource(), queueUrl)
		sqs.App().Flash().Info("Queue URL:" + queueUrl)
	}
	return nil
}
