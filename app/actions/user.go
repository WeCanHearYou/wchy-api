package actions

import (
	"context"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/i18n"
	"github.com/getfider/fider/app/pkg/validate"
)

//CreateUser is the action to create a new user
type CreateUser struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Reference string `json:"reference"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *CreateUser) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (action *CreateUser) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if action.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	} else if len(action.Name) > 100 {
		result.AddFieldFailure("name", "Name must have less than 100 characters.")
	}

	if action.Email == "" && action.Reference == "" {
		result.AddFieldFailure("", "Either email or reference is required")
	} else {
		if action.Email != "" {
			messages := validate.Email(ctx, action.Email)
			if len(messages) > 0 {
				result.AddFieldFailure("email", messages...)
			}
		}

		if len(action.Reference) > 100 {
			result.AddFieldFailure("reference", "Reference must have less than 100 characters.")
		}
	}

	return result
}

//ChangeUserRole is the input model change role of an user
type ChangeUserRole struct {
	Role   enum.Role `route:"role"`
	UserID int       `json:"userID"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *ChangeUserRole) IsAuthorized(ctx context.Context, user *entity.User) bool {
	if user == nil {
		return false
	}
	return user.IsAdministrator() && user.ID != action.UserID
}

// Validate if current model is valid
func (action *ChangeUserRole) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()
	if action.Role < enum.RoleVisitor || action.Role > enum.RoleAdministrator {
		return validate.Error(app.ErrNotFound)
	}

	if user.ID == action.UserID {
		result.AddFieldFailure("userID", "It is not allowed to change your own Role.")
	}

	userByID := &query.GetUserByID{UserID: action.UserID}
	err := bus.Dispatch(ctx, userByID)
	if err != nil {
		if errors.Cause(err) == app.ErrNotFound {
			result.AddFieldFailure("userID", "User not found.")
		} else {
			return validate.Error(err)
		}
	} else if userByID.Result.Tenant.ID != user.Tenant.ID {
		result.AddFieldFailure("userID", "User not found.")
	}
	return result
}

//ChangeUserEmail is the action used to change current user's email
type ChangeUserEmail struct {
	Email           string `json:"email" format:"lower"`
	VerificationKey string
	Requestor       *entity.User
}

func NewChangeUserEmail() *ChangeUserEmail {
	return &ChangeUserEmail{
		VerificationKey: entity.GenerateEmailVerificationKey(),
	}
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *ChangeUserEmail) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil
}

// Validate if current model is valid
func (action *ChangeUserEmail) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if action.Email == "" {
		result.AddFieldFailure("email", propertyIsRequired(ctx, "email"))
		return result
	}

	if user.Email == action.Email {
		result.AddFieldFailure("email", i18n.T(ctx, "validation.custom.differentemail"))
		return result
	}

	messages := validate.Email(ctx, action.Email)
	if len(messages) > 0 {
		result.AddFieldFailure("email", messages...)
		return result
	}

	userByEmail := &query.GetUserByEmail{Email: action.Email}
	err := bus.Dispatch(ctx, userByEmail)
	if err != nil && errors.Cause(err) != app.ErrNotFound {
		return validate.Error(err)
	}
	if err == nil && userByEmail.Result.ID != user.ID {
		result.AddFieldFailure("email", i18n.T(ctx, "validation.custom.emailtaken"))
		return result
	}
	action.Requestor = user
	return result
}

//GetEmail returns the email being verified
func (action *ChangeUserEmail) GetEmail() string {
	return action.Email
}

//GetName returns empty for this kind of process
func (action *ChangeUserEmail) GetName() string {
	return ""
}

//GetUser returns the current user performing this action
func (action *ChangeUserEmail) GetUser() *entity.User {
	return action.Requestor
}

//GetKind returns EmailVerificationKindSignIn
func (action *ChangeUserEmail) GetKind() enum.EmailVerificationKind {
	return enum.EmailVerificationKindChangeEmail
}
