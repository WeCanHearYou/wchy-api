package apiv1

import (
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/markdown"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/tasks"
)

// SendSampleInvite to current user's email
func SendSampleInvite() web.HandlerFunc {
	return func(c *web.Context) error {
		action := &actions.InviteUsers{IsSampleInvite: true}
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		if c.User().Email != "" {
			action.Message = strings.Replace(action.Message, app.InvitePlaceholder, "*the link to accept invitation will be here*", -1)
			to := dto.NewRecipient(c.User().Name, c.User().Email, dto.Props{
				"subject": action.Subject,
				"message": markdown.Full(action.Message),
			})

			bus.Publish(c, &cmd.SendMail{
				From:         c.Tenant().Name,
				To:           []dto.Recipient{to},
				TemplateName: "invite_email",
				Props: dto.Props{
					"logo": web.LogoURL(c),
				},
			})
		}

		return c.Ok(web.Map{})
	}
}

// SendInvites sends an email to each recipient
func SendInvites() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.InviteUsers)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		log.Warnf(c, "@{Tenant:magenta} sent @{TotalInvites:magenta} invites", dto.Props{
			"Tenant":       c.Tenant().Subdomain,
			"TotalInvites": len(action.Invitations),
		})
		c.Enqueue(tasks.SendInvites(action.Subject, action.Message, action.Invitations))

		return c.Ok(web.Map{})
	}
}
