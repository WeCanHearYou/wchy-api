import React, { useState } from "react"
import { Post } from "@fider/models"
import { Button, Icon } from "@fider/components"
import { actions } from "@fider/services"
import { useFider } from "@fider/hooks"
import IconVolumeOn from "@fider/assets/images/heroicons-volume-on.svg"
import IconVolumeOff from "@fider/assets/images/heroicons-volume-off.svg"
import { VStack } from "@fider/components/layout"
import { Trans } from "@lingui/macro"

interface NotificationsPanelProps {
  post: Post
  subscribed: boolean
}

export const NotificationsPanel = (props: NotificationsPanelProps) => {
  const fider = useFider()
  const [subscribed, setSubscribed] = useState(props.subscribed)

  const subscribeOrUnsubscribe = async () => {
    const action = subscribed ? actions.unsubscribe : actions.subscribe

    const response = await action(props.post.number)
    if (response.ok) {
      setSubscribed(!subscribed)
    }
  }

  if (!fider.session.isAuthenticated) {
    return null
  }

  const button = subscribed ? (
    <Button className="w-full" onClick={subscribeOrUnsubscribe}>
      <Icon sprite={IconVolumeOff} />{" "}
      <span>
        <Trans id="label.unsubscribe">Unsubscribe</Trans>
      </span>
    </Button>
  ) : (
    <Button className="w-full" onClick={subscribeOrUnsubscribe}>
      <Icon sprite={IconVolumeOn} />
      <span>
        <Trans id="label.subscribe">Subscribe</Trans>
      </span>
    </Button>
  )

  const text = subscribed ? (
    <span className="text-muted">
      <Trans id="showpost.notificationspanel.message.subscribed">You’re receiving notifications about activity on this post.</Trans>
    </span>
  ) : (
    <span className="text-muted">
      <Trans id="showpost.notificationspanel.message.unsubscribed">You&apos;ll not receive any notification about this post.</Trans>
    </span>
  )

  return (
    <VStack>
      <span className="text-category">
        <Trans id="label.notifications">Notifications</Trans>
      </span>
      {button}
      {text}
    </VStack>
  )
}
