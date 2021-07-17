import React, { useState } from "react"
import { SignInModal, TenantLogo, DevBanner, NotificationIndicator, UserMenu } from "@fider/components"
import { useFider } from "@fider/hooks"
import { HStack } from "./layout"
import { Trans } from "@lingui/macro"

export const Header = () => {
  const fider = useFider()
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)

  const showModal = (e: React.MouseEvent) => {
    e.preventDefault()
    setIsSignInModalOpen(true)
  }

  const hideModal = () => setIsSignInModalOpen(false)

  return (
    <>
      <DevBanner />
      <div id="c-header">
        <SignInModal isOpen={isSignInModalOpen} onClose={hideModal} />
        <HStack className="c-menu shadow p-2 w-full">
          <div className="container">
            <HStack justify="between">
              <a href="/" className="flex flex-x flex-items-center flex--spacing-2 h-8">
                <TenantLogo size={100} />
                <h1 className="text-title">{fider.session.tenant.name}</h1>
              </a>
              {fider.session.isAuthenticated && (
                <HStack spacing={2}>
                  <NotificationIndicator />
                  <UserMenu />
                </HStack>
              )}
              {!fider.session.isAuthenticated && (
                <a href="#" className="uppercase text-sm" onClick={showModal}>
                  <Trans id="action.signin">Sign in</Trans>
                </a>
              )}
            </HStack>
          </div>
        </HStack>
      </div>
    </>
  )
}
