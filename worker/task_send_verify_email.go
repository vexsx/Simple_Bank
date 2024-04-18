package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "github.com/vexsx/Simple-Bank/db/sqlc"
	"github.com/vexsx/Simple-Bank/util"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {

	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verfiy email : %w", err)
	}

	subject := "Welcome to Vexsx Simple Bank"
	// TODO: replace this URL with an environment variable that points to a front-end page
	verifyUrl := fmt.Sprintf("http://localhost:8080/Verify_Email?email_id=%d&secret_code=%s",
		verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(emailbody, nil, user.FullName, verifyUrl)
	to := []string{user.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).Msg("processed task")
	return nil
}

// email body design
const emailbody = `<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
        }
        .container {
            width: 80%;
            margin: auto;
            padding: 20px;
            border: 1px solid #ddd;
            border-radius: 5px;
        }
        .button {
            display: inline-block;
            padding: 10px 20px;
            margin: 10px 0px;
            border-radius: 5px;
            background-color: #4CAF50;
            color: white;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <h2>Hello %s,</h2>
        <p>Thank you for choosing to join us! We're excited to have you on board.</p>
        <p>As a final step to complete your registration, we kindly ask you to verify your email address by clicking on the button below:</p>
        <a href="%s" class="button">Verify Email Address</a>
        <p>If you didn't request this, please ignore this email. Your email won't be added to our list.</p>
        <p>Thank you for your time and welcome aboard!</p>
    </div>
</body>
</html>
	`
