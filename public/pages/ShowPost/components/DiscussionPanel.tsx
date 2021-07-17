import React from "react"
import { CurrentUser, Comment, Post } from "@fider/models"
import { ShowComment } from "./ShowComment"
import { CommentInput } from "./CommentInput"
import PostIllustration from "@fider/assets/images/undraw-post.svg"
import { Icon } from "@fider/components"
import { VStack } from "@fider/components/layout"
import { Trans } from "@lingui/macro"

interface DiscussionPanelProps {
  user?: CurrentUser
  post: Post
  comments: Comment[]
}

export const DiscussionPanel = (props: DiscussionPanelProps) => {
  return (
    <>
      <VStack spacing={2} className="c-comment-list">
        <span className="text-category">
          <Trans id="label.discussion">Discussion</Trans>
        </span>
        <VStack spacing={4} className="c-comment-list">
          {props.comments.map((c) => (
            <ShowComment key={c.id} post={props.post} comment={c} />
          ))}
          <CommentInput post={props.post} />
        </VStack>
        {props.comments.length === 0 && (
          <div className="text-center mt-4">
            <Icon sprite={PostIllustration} height="120" />
            <p className="text-xs">
              <Trans id="showpost.discussionpanel.emptymessage">No one has commented yet.</Trans>
            </p>
          </div>
        )}
      </VStack>
    </>
  )
}
