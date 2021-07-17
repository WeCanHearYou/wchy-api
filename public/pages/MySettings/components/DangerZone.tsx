import React from "react"

import { Button, Modal, ButtonClickEvent } from "@fider/components"
import { actions, notify, navigator } from "@fider/services"
import { Trans } from "@lingui/macro"

interface DangerZoneState {
  clicked: boolean
}

export class DangerZone extends React.Component<any, DangerZoneState> {
  constructor(props: any) {
    super(props)
    this.state = {
      clicked: false,
    }
  }

  public onClickDelete = async () => {
    this.setState({ clicked: true })
  }

  public onCancel = async () => {
    this.setState({ clicked: false })
  }

  public onConfirm = async (e: ButtonClickEvent) => {
    const response = await actions.deleteCurrentAccount()
    if (response.ok) {
      e.preventEnable()
      navigator.goHome()
    } else {
      notify.error("Failed to delete your account. Try again later")
    }
  }

  public render() {
    return (
      <div>
        <Modal.Window isOpen={this.state.clicked} center={false} onClose={this.onCancel}>
          <Modal.Header>
            <Trans id="modal.deleteaccount.header">Delete account</Trans>
          </Modal.Header>
          <Modal.Content>
            <Trans id="modal.deleteaccount.text">
              <p>
                When you choose to delete your account, we will erase all your personal information forever. The content you have published will remain, but it
                will be anonymised.
              </p>
              <p>
                This process is irreversible. <strong>Are you sure?</strong>
              </p>
            </Trans>
          </Modal.Content>
          <Modal.Footer>
            <Button variant="danger" size="small" onClick={this.onConfirm}>
              <Trans id="action.confirm">Confirm</Trans>
            </Button>
            <Button variant="tertiary" size="small" onClick={this.onCancel}>
              <Trans id="action.cancel">Cancel</Trans>
            </Button>
          </Modal.Footer>
        </Modal.Window>

        <h4 className="text-title mb-1">
          <Trans id="mysettings.dangerzone.title">Delete account</Trans>
        </h4>
        <p className="text-muted">
          <Trans id="mysettings.dangerzone.text">
            When you choose to delete your account, we will erase all your personal information forever. The content you have published will remain, but it will
            be anonymised.
          </Trans>
        </p>
        <p className="text-muted">
          <Trans id="mysettings.dangerzone.notice">This process is irreversible. Please be certain.</Trans>
        </p>
        <Button variant="danger" size="small" onClick={this.onClickDelete}>
          <Trans id="mysettings.dangerzone.delete">Delete My Account</Trans>
        </Button>
      </div>
    )
  }
}
