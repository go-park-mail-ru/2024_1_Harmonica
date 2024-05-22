package entity

import "html"

type Draft struct {
	SenderId   UserID `db:"sender_id" json:"sender_id"`
	ReceiverId UserID `db:"receiver_id" json:"receiver_id"`
	Text       string `db:"text" json:"text"`
}

func (d *Draft) Sanitize() {
	d.Text = html.EscapeString(d.Text)
}

type DraftResponse struct {
	Text string `db:"text" json:"text"`
}

func (d *DraftResponse) Sanitize() {
	d.Text = html.EscapeString(d.Text)
}
