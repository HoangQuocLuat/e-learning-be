package network

import "context"

func HasToken(ctx context.Context) bool {
	hasToken, _ := ctx.Value("has_token").(bool)
	return hasToken
}

func Token(ctx context.Context) string {
	token, _ := ctx.Value("token").(string)
	return token
}

func WorkspaceID(ctx context.Context) string {
	workspaceID, _ := ctx.Value("workspace_id").(string)
	return workspaceID
}

func SubWorkspaceID(ctx context.Context) string {
	workspaceID, _ := ctx.Value("sub_workspace_id").(string)
	return workspaceID
}

func AccountID(ctx context.Context) string {
	accountID, _ := ctx.Value("account_id").(string)
	return accountID
}

func Email(ctx context.Context) string {
	email, _ := ctx.Value("email").(string)
	return email
}
