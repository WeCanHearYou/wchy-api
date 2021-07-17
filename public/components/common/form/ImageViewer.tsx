import React from "react"
import { uploadedImageURL } from "@fider/services"
import { Modal, Button, Loader } from "@fider/components"

import "./ImageViewer.scss"
import { Trans } from "@lingui/macro"

interface ImageViewerProps {
  bkey: string
}

interface ImageViewerState {
  showModal: boolean
  loadedThumbnail: boolean
  loadedPreview: boolean
}

export class ImageViewer extends React.Component<ImageViewerProps, ImageViewerState> {
  constructor(props: ImageViewerProps) {
    super(props)

    this.state = {
      showModal: false,
      loadedThumbnail: false,
      loadedPreview: false,
    }
  }

  private openModal = () => {
    if (this.state.loadedThumbnail) {
      this.setState({ showModal: true })
    }
  }

  private closeModal = async () => {
    this.setState({ showModal: false })
  }

  private onThumbnailLoad = () => {
    this.setState({ loadedThumbnail: true })
  }

  private onPreviewLoad = () => {
    this.setState({ loadedPreview: true })
  }

  private modal() {
    return (
      <Modal.Window className="c-image-viewer-modal" isOpen={this.state.showModal} onClose={this.closeModal} center={false} size="fluid">
        <Modal.Content>
          {!this.state.loadedPreview && <Loader />}
          <img alt="" onLoad={this.onPreviewLoad} src={uploadedImageURL(this.props.bkey, 1500)} />
        </Modal.Content>

        <Modal.Footer>
          <Button variant="tertiary" onClick={this.closeModal}>
            <Trans id="action.close">Close</Trans>
          </Button>
        </Modal.Footer>
      </Modal.Window>
    )
  }

  public render() {
    const previewURL = uploadedImageURL(this.props.bkey, 200)
    return (
      <div className="c-image-viewer">
        {this.modal()}
        {!this.state.loadedThumbnail && <Loader />}
        <img alt="" onClick={this.openModal} onLoad={this.onThumbnailLoad} src={previewURL} />
      </div>
    )
  }
}
