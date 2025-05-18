package services

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/CustomCloudStorage/infrastructure/email"
	"github.com/CustomCloudStorage/repositories"
	"github.com/CustomCloudStorage/utils"
)

type EmailService interface {
	EnqueueEmail(ctx context.Context, tplName, to, subject string, data map[string]interface{}) error
	ProcessEmailQueue(ctx context.Context) error
}

type emailService struct {
	redis     repositories.RedisCache
	mailer    *email.SMTPMailer
	templates *template.Template
}

func NewEmailService(redis repositories.RedisCache, mailer *email.SMTPMailer, templates *template.Template) EmailService {
	svc := &emailService{
		redis:     redis,
		mailer:    mailer,
		templates: templates,
	}

	go func() {
		if err := svc.ProcessEmailQueue(context.Background()); err != nil {
			panic(err)
		}
	}()

	return svc
}

func (s *emailService) EnqueueEmail(ctx context.Context, tplName, to, subject string, data map[string]interface{}) error {
	if tplName == "" || to == "" {
		return utils.ErrBadRequest.Wrap(nil, "template name and recipient must be provided")
	}

	job := struct {
		Template string                 `json:"template"`
		To       string                 `json:"to"`
		Subject  string                 `json:"subject"`
		Data     map[string]interface{} `json:"data"`
	}{
		Template: tplName,
		To:       to,
		Subject:  subject,
		Data:     data,
	}

	err := s.redis.Enqueue(ctx, "email_queue", job)
	if err != nil {
		return utils.ErrInternal.Wrap(err, fmt.Sprintf("enqueue email to %s failed", to))
	}
	return nil
}

func (s *emailService) ProcessEmailQueue(ctx context.Context) error {
	for {
		var job struct {
			Template string                 `json:"template"`
			To       string                 `json:"to"`
			Subject  string                 `json:"subject"`
			Data     map[string]interface{} `json:"data"`
		}
		err := s.redis.Dequeue(ctx, "email_queue", &job)
		if err != nil {
			return utils.ErrInternal.Wrap(err, fmt.Sprintf("dequeue from %s failed", "email_queue"))
		}

		var buf bytes.Buffer
		err = s.templates.ExecuteTemplate(&buf, job.Template, job.Data)
		if err != nil {
			return utils.ErrInternal.Wrap(err, fmt.Sprintf("render template %q failed", job.Template))
		}

		err = s.mailer.Send(job.To, job.Subject, buf.String())
		if err != nil {
			return utils.ErrInternal.Wrap(err, fmt.Sprintf("send email to %s failed", job.To))
		}
	}
}
