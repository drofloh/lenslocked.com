package views

import "github.com/drofloh/lenslocked.com/models"

const (
	// AlertLvlError ...
	AlertLvlError = "danger"
	// AlertLvlWarning ...
	AlertLvlWarning = "warning"
	// AlertLvlInfo ...
	AlertLvlInfo = "info"
	// AlertLvlSuccess ...
	AlertLvlSuccess = "success"

	// AlertMsgGeneric is displayed when any random error is encountered by
	// our backend.
	AlertMsgGeneric = `Something went wrong. Please try again, and contact us
	if the problem persists.`
)

// Alert is used to render bootstrap alert messages in templates.
type Alert struct {
	Level   string
	Message string
}

// Data is the top level structure that views expect data to come in.
type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

// SetAlert ...
func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: pErr.Public(),
		}
	} else {
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: AlertMsgGeneric,
		}
	}
}

// AlertError ...
func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}

}

// PublicError ...
type PublicError interface {
	error
	Public() string
}
