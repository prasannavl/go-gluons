package responder

import (
	"fmt"
	"net/http"

	"html/template"

	"github.com/prasannavl/goerror/httperror"
	"github.com/unrolled/render"
)

// TODO: Proper content negotiation
// TODO: Use Content-Encoding

// Just use a singleton for now.
var Renderer = render.New(render.Options{})

func SetStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func SendErrorText(w http.ResponseWriter, errOrStringer interface{}) {
	var code int
	var message string
	switch e := errOrStringer.(type) {
	case error:
		message = e.Error()
		if e, ok := e.(httperror.HttpError); ok {
			code = e.Code()
		}
	case string:
		message = e
	case fmt.Stringer:
		message = e.String()
	}
	c := httperror.ErrorCode(code)
	if message == "" {
		SetStatus(w, c)
	} else {
		http.Error(w, message, c)
	}
}

func Send(w http.ResponseWriter, r *http.Request, value interface{}) error {
	return SendWithStatus(w, r, http.StatusOK, value)
}

func SendWithStatus(w http.ResponseWriter, r *http.Request, status int, value interface{}) error {
	// TODO: Do negotiation here
	if tcon, ok := value.(TemplateExecutor); ok {
		SetStatus(w, status)
		// TODO: Set proper headers before-hand. Meanwhile, net/http sniffs it, for now.
		return tcon.Execute(w)
	}
	return Renderer.JSON(w, status, value)
}

func SendError(w http.ResponseWriter, r *http.Request, err error) error {
	if e, ok := err.(httperror.HttpError); ok {
		return sendHttpError(w, r, e)
	}
	return SendWithStatus(w, r, http.StatusInternalServerError, err.Error())
}

func sendHttpError(w http.ResponseWriter, r *http.Request, err httperror.HttpError) error {
	msg := err.Error()
	return SendWithStatus(w, r, err.Code(), msg)
}

// Template stuff

type TemplateExecutor interface {
	Execute(w http.ResponseWriter) error
}

type TemplateFilesContent struct {
	Data          interface{}
	TemplateFiles []string
}

func (t *TemplateFilesContent) Execute(w http.ResponseWriter) error {
	tmp, err := template.ParseFiles(t.TemplateFiles...)
	if err != nil {
		return err
	}
	return tmp.Execute(w, t.Data)
}

type TemplateStringContent struct {
	Data           interface{}
	TemplateString string
}

func (t *TemplateStringContent) Execute(w http.ResponseWriter) error {
	tmp := template.New("tmpl")
	tmp, err := tmp.Parse(t.TemplateString)
	if err != nil {
		return err
	}
	return tmp.Execute(w, t.Data)
}

type TemplateContent struct {
	Data     interface{}
	Template template.Template
}

func (t *TemplateContent) Execute(w http.ResponseWriter) error {
	return t.Template.Execute(w, t.Data)
}

type TemplateGlobContent struct {
	Data interface{}
	Glob string
}

func (t *TemplateGlobContent) Execute(w http.ResponseWriter) error {
	tmp, err := template.ParseGlob(t.Glob)
	if err != nil {
		return err
	}
	return tmp.Execute(w, t.Data)
}
