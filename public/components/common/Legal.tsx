import React from "react"
import { Modal, Checkbox } from "@fider/components"
import { useFider } from "@fider/hooks"
import { Trans } from "@lingui/macro"

interface LegalAgreementProps {
  onChange: (agreed: boolean) => void
}

export const TermsOfService = () => {
  const fider = useFider()

  if (fider.settings.hasLegal) {
    return (
      <a href="/terms" className="text-link" target="_blank">
        <Trans id="legal.termsofservice">Terms of Service</Trans>
      </a>
    )
  }
  return null
}

export const PrivacyPolicy = () => {
  const fider = useFider()

  if (fider.settings.hasLegal) {
    return (
      <a href="/privacy" className="text-link" target="_blank">
        <Trans id="legal.privacypolicy">Privacy Policy</Trans>
      </a>
    )
  }
  return null
}

export const LegalNotice = () => {
  const fider = useFider()

  if (fider.settings.hasLegal) {
    return (
      <p className="text-muted">
        <Trans id="legal.notice">
          By signing in, you agree to the <PrivacyPolicy /> and <TermsOfService />.
        </Trans>
      </p>
    )
  }
  return null
}

export const LegalFooter = () => {
  const fider = useFider()

  if (fider.settings.hasLegal) {
    return (
      <Modal.Footer align="center">
        <LegalNotice />
      </Modal.Footer>
    )
  }
  return null
}

export const LegalAgreement: React.FunctionComponent<LegalAgreementProps> = (props) => {
  const fider = useFider()

  if (fider.settings.hasLegal) {
    return (
      <Checkbox field="legalAgreement" onChange={props.onChange}>
        <Trans id="legal.agreement">
          I have read and agree to the <PrivacyPolicy /> and <TermsOfService />.
        </Trans>
      </Checkbox>
    )
  }
  return null
}
