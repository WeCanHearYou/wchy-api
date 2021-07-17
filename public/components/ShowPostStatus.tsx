import "./ShowPostStatus.scss"

import React from "react"
import { PostStatus } from "@fider/models"
import { t } from "@lingui/macro"

interface ShowPostStatusProps {
  status: PostStatus
}

export const ShowPostStatus = (props: ShowPostStatusProps) => {
  const id = `enum.poststatus.${props.status.value}`
  const title = t({ id, message: props.status.title })
  return <span className={`c-status-label c-status-label--${props.status.value}`}>{title}</span>
}
