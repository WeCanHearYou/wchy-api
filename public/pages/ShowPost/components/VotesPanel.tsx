import React, { useState } from "react"
import { Post, Vote } from "@fider/models"
import { AvatarStack, Button } from "@fider/components"
import { Fider } from "@fider/services"
import { useFider } from "@fider/hooks"
import { VotesModal } from "./VotesModal"
import { HStack, VStack } from "@fider/components/layout"
import { Trans } from "@lingui/macro"

interface VotesPanelProps {
  post: Post
  votes: Vote[]
}

export const VotesPanel = (props: VotesPanelProps) => {
  const fider = useFider()
  const [isVotesModalOpen, setIsVotesModalOpen] = useState(false)
  const canShowAll = fider.session.isAuthenticated && Fider.session.user.isCollaborator

  const openModal = () => {
    if (canShowAll) {
      setIsVotesModalOpen(true)
    }
  }

  const closeModal = () => setIsVotesModalOpen(false)

  const extraVotesCount = props.post.votesCount - props.votes.length

  return (
    <VStack>
      <VotesModal post={props.post} isOpen={isVotesModalOpen} onClose={closeModal} />
      <span className="text-category">
        <Trans id="label.voters">Voters</Trans>
      </span>
      <HStack>
        {props.votes.length > 0 && <AvatarStack users={props.votes.map((x) => x.user)} />}
        {extraVotesCount > 0 && (
          <Button variant="tertiary" disabled={!canShowAll} size="small" onClick={openModal}>
            <Trans id="showpost.votespanel.more">+{extraVotesCount} more</Trans>
          </Button>
        )}
        {props.votes.length > 0 && extraVotesCount === 0 && canShowAll && (
          <Button variant="tertiary" size="small" disabled={!canShowAll} onClick={openModal}>
            <Trans id="showpost.votespanel.seedetails">see details</Trans>
          </Button>
        )}
        {props.votes.length === 0 && (
          <span className="text-muted">
            <Trans id="label.none">None</Trans>
          </span>
        )}
      </HStack>
    </VStack>
  )
}
