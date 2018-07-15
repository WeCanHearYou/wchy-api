import * as React from "react";
import { Button } from "@fider/components/common";
import { OAuthProviderOption } from "@fider/models";
import { classSet } from "@fider/services";

interface OAuthProviderPreview {
  provider: string;
  displayName: string;
}

interface SocialSignInButtonProps {
  option: OAuthProviderOption | OAuthProviderPreview;
  redirectTo?: string;
}

export class SocialSignInButton extends React.Component<SocialSignInButtonProps, {}> {
  public render() {
    const redirectTo = this.props.redirectTo || location.href;
    const href = "url" in this.props.option ? `${this.props.option.url}?redirect=${redirectTo}` : undefined;
    const className = classSet({
      "m-social": true,
      [`m-${this.props.option.provider}`]: true
    });

    return (
      <Button href={href} fluid={true} className={className}>
        <i className="svg" />
        <span>{this.props.option.displayName}</span>
      </Button>
    );
  }
}
