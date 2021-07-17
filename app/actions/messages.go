package actions

import (
	"context"
	"fmt"

	"github.com/getfider/fider/app/pkg/i18n"
)

func propertyIsRequired(ctx context.Context, fieldName string) string {
	displayName := i18n.T(ctx, fmt.Sprintf("property.%s", fieldName))
	return i18n.T(ctx, "validation.required", i18n.Params{"name": displayName})
}

func propertyIsInvalid(ctx context.Context, fieldName string) string {
	displayName := i18n.T(ctx, fmt.Sprintf("property.%s", fieldName))
	return i18n.T(ctx, "validation.invalid", i18n.Params{"name": displayName})
}

func propertyMaxStringLen(ctx context.Context, fieldName string, maxLen int) string {
	displayName := i18n.T(ctx, fmt.Sprintf("property.%s", fieldName))
	return i18n.T(ctx, "validation.maxstringlen",
		i18n.Params{"name": displayName},
		i18n.Params{"len": maxLen},
	)
}
