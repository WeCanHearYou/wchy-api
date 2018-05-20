import * as React from "react";
import { IdeaStatus, CurrentUser, Idea } from "@fider/models";
import { page, actions, Failure } from "@fider/services";
import { Form, DisplayError, Textarea, Modal, Button, List, ListItem, Form2, TextArea } from "@fider/components";

interface ModerationPanelProps {
  user?: CurrentUser;
  idea: Idea;
}

interface ModerationPanelState {
  showConfirmation: boolean;
  text: string;
  error?: Failure;
}

export class ModerationPanel extends React.Component<ModerationPanelProps, ModerationPanelState> {
  private form!: Form;

  constructor(props: ModerationPanelProps) {
    super(props);
    this.state = {
      text: "",
      showConfirmation: false
    };
  }

  private async delete(): Promise<void> {
    const response = await actions.deleteIdea(this.props.idea.number, this.state.text);
    if (response.ok) {
      this.close();
      page.goHome();
    } else if (response.error) {
      this.setState({ error: this.state.error });
    }
  }

  private async close(): Promise<void> {
    this.setState({ showConfirmation: false });
  }

  public render() {
    const status = IdeaStatus.Get(this.props.idea.status);
    if (!this.props.user || !this.props.user.isAdministrator || status.closed) {
      return null;
    }

    const modal = (
      <Modal.Window isOpen={this.state.showConfirmation} center={false} size="large">
        <Modal.Content>
          <Form2 error={this.state.error}>
            <TextArea
              field="text"
              onChange={text => this.setState({ text })}
              value={this.state.text}
              placeholder="Why are you deleting this idea? (optional)"
            >
              <span className="info">
                This operation <strong>cannot</strong> be undone.
              </span>
            </TextArea>
          </Form2>
        </Modal.Content>

        <Modal.Footer>
          <Button color="danger" onClick={async () => this.delete()}>
            Delete
          </Button>
          <Button onClick={async () => this.close()}>Cancel</Button>
        </Modal.Footer>
      </Modal.Window>
    );

    return (
      <div>
        {modal}
        <span className="subtitle">Moderation</span>
        <List>
          <ListItem>
            <Button
              color="danger"
              size="tiny"
              fluid={true}
              onClick={async () => this.setState({ showConfirmation: true })}
            >
              <i className="delete icon" /> Delete
            </Button>
          </ListItem>
        </List>
      </div>
    );
  }
}
