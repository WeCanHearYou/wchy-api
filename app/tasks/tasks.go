package tasks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/i18n"
	"github.com/getfider/fider/app/pkg/markdown"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/worker"
)

func describe(name string, job worker.Job) worker.Task {
	return worker.Task{Name: name, Job: job}
}

func link(baseURL, path string, args ...interface{}) string {
	return fmt.Sprintf("<a href='%[1]s%[2]s'>%[1]s%[2]s</a>", baseURL, fmt.Sprintf(path, args...))
}

func linkWithText(text, baseURL, path string, args ...interface{}) string {
	return fmt.Sprintf("<a href='%s%s'>%s</a>", baseURL, fmt.Sprintf(path, args...), text)
}

//SendSignUpEmail is used to send the sign up email to requestor
func SendSignUpEmail(action *actions.CreateTenant, baseURL string) worker.Task {
	return describe("Send sign up email", func(c *worker.Context) error {
		to := dto.NewRecipient(action.Name, action.Email, dto.Props{
			"link": link(baseURL, "/signup/verify?k=%s", action.VerificationKey),
		})

		bus.Publish(c, &cmd.SendMail{
			From:         "Fider",
			To:           []dto.Recipient{to},
			TemplateName: "signup_email",
			Props: dto.Props{
				"logo": web.LogoURL(c),
			},
		})

		return nil
	})
}

//SendSignInEmail is used to send the sign in email to requestor
func SendSignInEmail(email, verificationKey string) worker.Task {
	return describe("Send sign in email", func(c *worker.Context) error {
		to := dto.NewRecipient("", email, dto.Props{
			"siteName": c.Tenant().Name,
			"link":     link(web.BaseURL(c), "/signin/verify?k=%s", verificationKey),
		})

		bus.Publish(c, &cmd.SendMail{
			From:         c.Tenant().Name,
			To:           []dto.Recipient{to},
			TemplateName: "signin_email",
			Props: dto.Props{
				"logo": web.LogoURL(c),
			},
		})

		return nil
	})
}

//SendChangeEmailConfirmation is used to send the change email confirmation email to requestor
func SendChangeEmailConfirmation(action *actions.ChangeUserEmail) worker.Task {
	return describe("Send change email confirmation", func(c *worker.Context) error {

		previous := c.User().Email
		if previous == "" {
			previous = "(empty)"
		}

		to := dto.NewRecipient(action.Requestor.Name, action.Email, dto.Props{
			"name":     c.User().Name,
			"oldEmail": previous,
			"newEmail": action.Email,
			"link":     link(web.BaseURL(c), "/change-email/verify?k=%s", action.VerificationKey),
		})

		bus.Publish(c, &cmd.SendMail{
			From:         c.Tenant().Name,
			To:           []dto.Recipient{to},
			TemplateName: "change_emailaddress_email",
			Props: dto.Props{
				"logo": web.LogoURL(c),
			},
		})

		return nil
	})
}

//NotifyAboutNewPost sends a notification (web and email) to subscribers
func NotifyAboutNewPost(post *entity.Post) worker.Task {
	return describe("Notify about new post", func(c *worker.Context) error {
		// Web notification
		users, err := getActiveSubscribers(c, post, enum.NotificationChannelWeb, enum.NotificationEventNewPost)
		if err != nil {
			return c.Failure(err)
		}

		title := fmt.Sprintf("New post: **%s**", post.Title)
		link := fmt.Sprintf("/posts/%d/%s", post.Number, post.Slug)
		for _, user := range users {
			if user.ID != c.User().ID {
				err = bus.Dispatch(c, &cmd.AddNewNotification{
					User:   user,
					Title:  title,
					Link:   link,
					PostID: post.ID,
				})
				if err != nil {
					return c.Failure(err)
				}
			}
		}

		// Email notification
		users, err = getActiveSubscribers(c, post, enum.NotificationChannelEmail, enum.NotificationEventNewPost)
		if err != nil {
			return c.Failure(err)
		}

		to := make([]dto.Recipient, 0)
		for _, user := range users {
			if user.ID != c.User().ID {
				to = append(to, dto.NewRecipient(user.Name, user.Email, dto.Props{}))
			}
		}

		props := dto.Props{
			"title":    post.Title,
			"siteName": c.Tenant().Name,
			"userName": c.User().Name,
			"content":  markdown.Full(post.Description),
			"postLink": linkWithText(fmt.Sprintf("#%d", post.Number), web.BaseURL(c), "/posts/%d/%s", post.Number, post.Slug),
			"view":     linkWithText(i18n.T(c, "email.subscription.view"), web.BaseURL(c), "/posts/%d/%s", post.Number, post.Slug),
			"change":   linkWithText(i18n.T(c, "email.subscription.change"), web.BaseURL(c), "/settings"),
			"logo":     web.LogoURL(c),
		}

		bus.Publish(c, &cmd.SendMail{
			From:         c.User().Name,
			To:           to,
			TemplateName: "new_post",
			Props:        props,
		})

		return nil
	})
}

//NotifyAboutNewComment sends a notification (web and email) to subscribers
func NotifyAboutNewComment(post *entity.Post, comment string) worker.Task {
	return describe("Notify about new comment", func(c *worker.Context) error {
		// Web notification
		users, err := getActiveSubscribers(c, post, enum.NotificationChannelWeb, enum.NotificationEventNewComment)
		if err != nil {
			return c.Failure(err)
		}

		title := fmt.Sprintf("**%s** left a comment on **%s**", c.User().Name, post.Title)
		link := fmt.Sprintf("/posts/%d/%s", post.Number, post.Slug)
		for _, user := range users {
			if user.ID != c.User().ID {
				err = bus.Dispatch(c, &cmd.AddNewNotification{
					User:   user,
					Title:  title,
					Link:   link,
					PostID: post.ID,
				})
				if err != nil {
					return c.Failure(err)
				}
			}
		}

		// Email notification
		users, err = getActiveSubscribers(c, post, enum.NotificationChannelEmail, enum.NotificationEventNewComment)
		if err != nil {
			return c.Failure(err)
		}

		to := make([]dto.Recipient, 0)
		for _, user := range users {
			if user.ID != c.User().ID {
				to = append(to, dto.NewRecipient(user.Name, user.Email, dto.Props{}))
			}
		}

		props := dto.Props{
			"title":       post.Title,
			"siteName":    c.Tenant().Name,
			"userName":    c.User().Name,
			"content":     markdown.Full(comment),
			"postLink":    linkWithText(fmt.Sprintf("#%d", post.Number), web.BaseURL(c), "/posts/%d/%s", post.Number, post.Slug),
			"view":        linkWithText(i18n.T(c, "email.subscription.view"), web.BaseURL(c), "/posts/%d/%s", post.Number, post.Slug),
			"unsubscribe": linkWithText(i18n.T(c, "email.subscription.unsubscribe"), web.BaseURL(c), "/posts/%d/%s", post.Number, post.Slug),
			"change":      linkWithText(i18n.T(c, "email.subscription.change"), web.BaseURL(c), "/settings"),
			"logo":        web.LogoURL(c),
		}

		bus.Publish(c, &cmd.SendMail{
			From:         c.User().Name,
			To:           to,
			TemplateName: "new_comment",
			Props:        props,
		})

		return nil
	})
}

//NotifyAboutStatusChange sends a notification (web and email) to subscribers
func NotifyAboutStatusChange(post *entity.Post, prevStatus enum.PostStatus) worker.Task {
	return describe("Notify about post status change", func(c *worker.Context) error {
		//Don't notify if previous status is the same
		if prevStatus == post.Status {
			return nil
		}

		// Web notification
		users, err := getActiveSubscribers(c, post, enum.NotificationChannelWeb, enum.NotificationEventChangeStatus)
		if err != nil {
			return c.Failure(err)
		}

		title := fmt.Sprintf("**%s** changed status of **%s** to **%s**", c.User().Name, post.Title, post.Status.Name())
		link := fmt.Sprintf("/posts/%d/%s", post.Number, post.Slug)
		for _, user := range users {
			if user.ID != c.User().ID {
				err = bus.Dispatch(c, &cmd.AddNewNotification{
					User:   user,
					Title:  title,
					Link:   link,
					PostID: post.ID,
				})
				if err != nil {
					return c.Failure(err)
				}
			}
		}

		// Email notification
		users, err = getActiveSubscribers(c, post, enum.NotificationChannelEmail, enum.NotificationEventChangeStatus)
		if err != nil {
			return c.Failure(err)
		}

		var duplicate string
		if post.Status == enum.PostDuplicate {
			duplicate = linkWithText(post.Response.Original.Title, web.BaseURL(c), "/posts/%d/%s", post.Response.Original.Number, post.Response.Original.Slug)
		}

		to := make([]dto.Recipient, 0)
		for _, user := range users {
			if user.ID != c.User().ID {
				to = append(to, dto.NewRecipient(user.Name, user.Email, dto.Props{}))
			}
		}

		props := dto.Props{
			"title":       post.Title,
			"postLink":    linkWithText(fmt.Sprintf("#%d", post.Number), web.BaseURL(c), "/posts/%d/%s", post.Number, post.Slug),
			"siteName":    c.Tenant().Name,
			"content":     markdown.Full(post.Response.Text),
			"status":      i18n.T(c, fmt.Sprintf("enum.poststatus.%s", post.Status.Name())),
			"duplicate":   duplicate,
			"view":        linkWithText(i18n.T(c, "email.subscription.view"), web.BaseURL(c), "/posts/%d/%s", post.Number, post.Slug),
			"unsubscribe": linkWithText(i18n.T(c, "email.subscription.unsubscribe"), web.BaseURL(c), "/posts/%d/%s", post.Number, post.Slug),
			"change":      linkWithText(i18n.T(c, "email.subscription.change"), web.BaseURL(c), "/settings"),
			"logo":        web.LogoURL(c),
		}

		bus.Publish(c, &cmd.SendMail{
			From:         c.User().Name,
			To:           to,
			TemplateName: "change_status",
			Props:        props,
		})

		return nil
	})
}

//NotifyAboutDeletedPost sends a notification (web and email) to subscribers of the post that has been deleted
func NotifyAboutDeletedPost(post *entity.Post) worker.Task {
	return describe("Notify about deleted post", func(c *worker.Context) error {

		// Web notification
		users, err := getActiveSubscribers(c, post, enum.NotificationChannelWeb, enum.NotificationEventChangeStatus)
		if err != nil {
			return c.Failure(err)
		}

		title := fmt.Sprintf("**%s** deleted **%s**", c.User().Name, post.Title)
		for _, user := range users {
			if user.ID != c.User().ID {
				err = bus.Dispatch(c, &cmd.AddNewNotification{
					User:   user,
					Title:  title,
					PostID: post.ID,
				})
				if err != nil {
					return c.Failure(err)
				}
			}
		}

		// Email notification
		users, err = getActiveSubscribers(c, post, enum.NotificationChannelEmail, enum.NotificationEventChangeStatus)
		if err != nil {
			return c.Failure(err)
		}

		to := make([]dto.Recipient, 0)
		for _, user := range users {
			if user.ID != c.User().ID {
				to = append(to, dto.NewRecipient(user.Name, user.Email, dto.Props{}))
			}
		}

		props := dto.Props{
			"title":    post.Title,
			"siteName": c.Tenant().Name,
			"content":  markdown.Full(post.Response.Text),
			"change":   linkWithText(i18n.T(c, "email.subscription.change"), web.BaseURL(c), "/settings"),
			"logo":     web.LogoURL(c),
		}

		bus.Publish(c, &cmd.SendMail{
			From:         c.User().Name,
			To:           to,
			TemplateName: "delete_post",
			Props:        props,
		})

		return nil
	})
}

//SendInvites sends one email to each invited recipient
func SendInvites(subject, message string, invitations []*actions.UserInvitation) worker.Task {
	return describe("Send invites", func(c *worker.Context) error {
		to := make([]dto.Recipient, len(invitations))
		for i, invite := range invitations {
			err := bus.Dispatch(c, &cmd.SaveVerificationKey{
				Key:      invite.VerificationKey,
				Duration: 15 * 24 * time.Hour,
				Request:  invite,
			})
			if err != nil {
				return c.Failure(err)
			}

			url := fmt.Sprintf("%s/invite/verify?k=%s", web.BaseURL(c), invite.VerificationKey)
			toMessage := strings.Replace(message, app.InvitePlaceholder, string(url), -1)
			to[i] = dto.NewRecipient("", invite.Email, dto.Props{
				"message": markdown.Full(toMessage),
			})
		}

		bus.Publish(c, &cmd.SendMail{
			From:         c.User().Name,
			To:           to,
			TemplateName: "invite_email",
			Props: dto.Props{
				"subject": subject,
				"logo":    web.LogoURL(c),
			},
		})

		return nil
	})
}

func getActiveSubscribers(ctx context.Context, post *entity.Post, channel enum.NotificationChannel, event enum.NotificationEvent) ([]*entity.User, error) {
	q := &query.GetActiveSubscribers{
		Number:  post.Number,
		Channel: channel,
		Event:   event,
	}
	err := bus.Dispatch(ctx, q)
	return q.Result, err
}
