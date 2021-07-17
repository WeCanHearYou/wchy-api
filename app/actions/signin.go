package actions

import (
	"context"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/validate"
)

// SignInByEmail happens when user request to sign in by email
type SignInByEmail struct {
	Email           string `json:"email" format:"lower"`
	VerificationKey string
}

func NewSignInByEmail() *SignInByEmail {
	return &SignInByEmail{
		VerificationKey: entity.GenerateEmailVerificationKey(),
	}
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *SignInByEmail) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return true
}

// Validate if current model is valid
func (action *SignInByEmail) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if action.Email == "" {
		result.AddFieldFailure("email", propertyIsRequired(ctx, "email"))
		return result
	}

	messages := validate.Email(ctx, action.Email)
	result.AddFieldFailure("email", messages...)

	return result
}

//GetEmail returns the email being verified
func (action *SignInByEmail) GetEmail() string {
	return action.Email
}

//GetName returns empty for this kind of process
func (action *SignInByEmail) GetName() string {
	return ""
}

//GetUser returns the current user performing this action
func (action *SignInByEmail) GetUser() *entity.User {
	return nil
}

//GetKind returns EmailVerificationKindSignIn
func (action *SignInByEmail) GetKind() enum.EmailVerificationKind {
	return enum.EmailVerificationKindSignIn
}

// CompleteProfile happens when users completes their profile during first time sign in
type CompleteProfile struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Email string
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *CompleteProfile) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return true
}

// Validate if current model is valid
func (action *CompleteProfile) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if action.Name == "" {
		result.AddFieldFailure("name", propertyIsRequired(ctx, "name"))
	} else if len(action.Name) > 50 {
		result.AddFieldFailure("name", propertyMaxStringLen(ctx, "name", 50))
	}

	if action.Key == "" {
		result.AddFieldFailure("key", propertyIsRequired(ctx, "key"))
	} else {
		findBySignIn := &query.GetVerificationByKey{Kind: enum.EmailVerificationKindSignIn, Key: action.Key}
		err1 := bus.Dispatch(ctx, findBySignIn)

		findByUserInvitation := &query.GetVerificationByKey{Kind: enum.EmailVerificationKindUserInvitation, Key: action.Key}
		err2 := bus.Dispatch(ctx, findByUserInvitation)

		if err1 == nil {
			action.Email = findBySignIn.Result.Email
		} else if err2 == nil {
			action.Email = findByUserInvitation.Result.Email
		} else {
			result.AddFieldFailure("key", propertyIsInvalid(ctx, "key"))
		}
	}

	return result
}
