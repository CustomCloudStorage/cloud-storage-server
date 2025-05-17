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

type EmailPayload struct {
	Template string                 `json:"template"`
	To       string                 `json:"to"`
	Subject  string                 `json:"subject"`
	Data     map[string]interface{} `json:"data"`
}

type EmailService struct {
	queueName string
	redis     repositories.RedisCache
	mailer    *email.SMTPMailer
	templates *template.Template
}

func NewEmailService(queueName string, redis repositories.RedisCache, mailer *email.SMTPMailer, templates *template.Template) *EmailService {
	svc := &EmailService{
		queueName: queueName,
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

func (s *EmailService) EnqueueEmail(ctx context.Context, tplName, to, subject string, data map[string]interface{}) error {
	if tplName == "" || to == "" {
		return utils.ErrBadRequest.Wrap(nil, "template name and recipient must be provided")
	}

	job := EmailPayload{
		Template: tplName,
		To:       to,
		Subject:  subject,
		Data:     data,
	}

	if err := s.redis.Enqueue(ctx, s.queueName, job); err != nil {
		return utils.ErrInternal.Wrap(err, fmt.Sprintf("enqueue email to %s failed", to))
	}
	return nil
}

func (s *EmailService) ProcessEmailQueue(ctx context.Context) error {
	for {
		var job EmailPayload
		if err := s.redis.Dequeue(ctx, s.queueName, &job); err != nil {
			return utils.ErrInternal.Wrap(err, fmt.Sprintf("dequeue from %s failed", s.queueName))
		}

		var buf bytes.Buffer
		if err := s.templates.ExecuteTemplate(&buf, job.Template, job.Data); err != nil {
			return utils.ErrInternal.Wrap(err, fmt.Sprintf("render template %q failed", job.Template))
		}

		if err := s.mailer.Send(job.To, job.Subject, buf.String()); err != nil {
			return utils.ErrInternal.Wrap(err, fmt.Sprintf("send email to %s failed", job.To))
		}
	}
}
