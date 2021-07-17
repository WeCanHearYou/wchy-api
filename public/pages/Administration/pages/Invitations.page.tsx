import React from "react"

import { Button, TextArea, Form, Input, Field } from "@fider/components"
import { actions, notify, Failure, Fider } from "@fider/services"
import { AdminBasePage } from "../components/AdminBasePage"

interface InvitationsPageState {
  subject: string
  message: string
  recipients: string[]
  numOfRecipients: number
  rawRecipients: string
  error?: Failure
}

export default class InvitationsPage extends AdminBasePage<any, InvitationsPageState> {
  public id = "p-admin-invitations"
  public name = "invitations"
  public title = "Invitations"
  public subtitle = "Invite people to share their feedback"

  constructor(props: any) {
    super(props)

    this.state = {
      subject: `Share your ideas and thoughts about ${Fider.session.tenant.name}`,
      message: `Hi,

At **${Fider.session.tenant.name}** we take feedback very seriously, which is why we've launched a space where you can vote, discuss and share your ideas and thoughts about our products and services.

We'd like to extend an invite for you to join this community and raise awareness on topics you care about!

To join, click on the link below.

%invite%

Regards,
${Fider.session.user.name} (${Fider.session.tenant.name})`,
      recipients: [],
      numOfRecipients: 0,
      rawRecipients: "",
    }
  }

  private setRecipients = (rawRecipients: string) => {
    const recipients = rawRecipients.split(/\n|;|,|\s/gm).filter((x) => !!x)
    this.setState({ rawRecipients, recipients, numOfRecipients: recipients.length })
  }

  private sendSample = async () => {
    const result = await actions.sendSampleInvite(this.state.subject, this.state.message)
    if (result.ok) {
      notify.success(
        <span>
          An email message was sent to <strong>{Fider.session.user.email}</strong>
        </span>
      )
    }
    this.setState({ error: result.error })
  }

  private sendInvites = async () => {
    const result = await actions.sendInvites(this.state.subject, this.state.message, this.state.recipients)
    if (result.ok) {
      notify.success("Your invites have been sent.")
      this.setState({ rawRecipients: "", numOfRecipients: 0, recipients: [], error: undefined })
    } else {
      this.setState({ error: result.error })
    }
  }

  private setSubject = (subject: string): void => {
    this.setState({ subject })
  }

  private setMessage = (message: string): void => {
    this.setState({ message })
  }

  public content() {
    return (
      <Form error={this.state.error}>
        <TextArea
          field="recipients"
          label="Send invitations to"
          placeholder="james@example.com; mary@example.com"
          minRows={1}
          value={this.state.rawRecipients}
          onChange={this.setRecipients}
        >
          <div className="text-muted">
            <p>
              Input the list of all email addresses you wish to invite. Separate each address with either <strong>semicolon</strong>, <strong>comma</strong>,{" "}
              <strong>whitespace</strong> or <strong>line break</strong>.
            </p>
            <p>You can send this invite to a maximum of 30 recipients each time.</p>
          </div>
        </TextArea>

        <Input field="subject" label="Subject" value={this.state.subject} maxLength={70} onChange={this.setSubject}>
          <p className="text-muted">This is the subject that will be used on the invitation email. Keep it short and sweet.</p>
        </Input>

        <TextArea field="message" label="Message" minRows={8} value={this.state.message} onChange={this.setMessage}>
          <div className="text-muted">
            <p>
              This is the content of the invite. Be polite and explain what this invite is for, otherwise there&apos;s a high change people will ignore your
              message.
            </p>
            <p>
              You&apos;re allowed to write whatever you want as long as you include the invitation link placeholder named <strong>%invite%</strong>.
            </p>
          </div>
        </TextArea>

        <Field label="Sample Invite">
          <p className="text-muted">
            We highly recommend to send yourself a sample email for you to verify if everything is correct before inviting your list of contacts.
          </p>
          {Fider.session.user.email ? (
            <Button onClick={this.sendSample}>Send a sample email to {Fider.session.user.email}</Button>
          ) : (
            <Button disabled={true}>Your profile doesn&apos;t have an email</Button>
          )}
        </Field>

        <Field label="Confirmation">
          <p className="text-muted">Whenever you&apos;re ready, click the following button to send out these invites.</p>
          <Button onClick={this.sendInvites} variant="primary" disabled={this.state.numOfRecipients === 0}>
            Send {this.state.numOfRecipients} {this.state.numOfRecipients === 1 ? "invite" : "invites"}
          </Button>
        </Field>
      </Form>
    )
  }
}
