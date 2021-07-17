import React from "react"
import { Post, Vote } from "@fider/models"
import { Modal, Button, Loader, Avatar, UserName, Moment, Input } from "@fider/components"
import { actions, Fider } from "@fider/services"
import IconSearch from "@fider/assets/images/heroicons-search.svg"
import IconX from "@fider/assets/images/heroicons-x.svg"
import { HStack, VStack } from "@fider/components/layout"
import { t, Trans } from "@lingui/macro"

interface VotesModalProps {
  isOpen: boolean
  post: Post
  onClose?: () => void
}

interface VotesModalState {
  searchText: string
  allVotes: Vote[]
  filteredVotes: Vote[]
  isLoading: boolean
  query: string
}

export class VotesModal extends React.Component<VotesModalProps, VotesModalState> {
  constructor(props: VotesModalProps) {
    super(props)
    this.state = {
      searchText: "",
      query: "",
      allVotes: [],
      filteredVotes: [],
      isLoading: true,
    }
  }

  public componentDidUpdate(prevProps: VotesModalProps) {
    if (this.props.isOpen && !prevProps.isOpen) {
      actions.listVotes(this.props.post.number).then((response) => {
        if (response.ok) {
          this.setState({
            allVotes: response.data,
            filteredVotes: response.data,
            isLoading: false,
          })
        }
      })
    }
  }

  private closeModal = async () => {
    if (this.props.onClose) {
      this.props.onClose()
    }
  }

  private clearSearch = () => {
    this.handleSearchFilterChanged("")
  }

  private handleSearchFilterChanged = (query: string) => {
    const votes = this.state.allVotes.filter((x) => x.user.name.toLowerCase().indexOf(query.toLowerCase()) >= 0)
    this.setState({ query, filteredVotes: votes })
  }

  public render() {
    return (
      <Modal.Window isOpen={this.props.isOpen} center={false} onClose={this.closeModal}>
        <Modal.Content>
          {this.state.isLoading && <Loader />}
          {!this.state.isLoading && (
            <>
              <Input
                field="query"
                icon={this.state.query ? IconX : IconSearch}
                onIconClick={this.state.query ? this.clearSearch : undefined}
                placeholder={t({ id: "modal.showvotes.query.placeholder", message: "Search for users by name..." })}
                value={this.state.query}
                onChange={this.handleSearchFilterChanged}
              />
              <VStack spacing={2} className="h-max-5xl overflow-scroll">
                {this.state.filteredVotes.map((x) => (
                  <HStack key={x.user.id} justify="between">
                    <HStack>
                      <Avatar user={x.user} />
                      <VStack spacing={0}>
                        <UserName user={x.user} />
                        <span className="text-muted">{x.user.email}</span>
                      </VStack>
                    </HStack>
                    <span className="text-muted">
                      <Moment locale={Fider.currentLocale} date={x.createdAt} />
                    </span>
                  </HStack>
                ))}
                {this.state.filteredVotes.length === 0 && (
                  <p className="text-muted">
                    <Trans id="modal.showvotes.message.zeromatches">No users found matching &apos;{this.state.query}&apos;.</Trans>
                  </p>
                )}
              </VStack>
            </>
          )}
        </Modal.Content>

        <Modal.Footer>
          <Button variant="tertiary" onClick={this.closeModal}>
            <Trans id="action.close">Close</Trans>
          </Button>
        </Modal.Footer>
      </Modal.Window>
    )
  }
}
