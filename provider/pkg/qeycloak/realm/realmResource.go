package realm

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/raushan606/pulumi-qeycloak/provider/pkg/qeycloak/config"
)

type Realm struct{}

type RealmArgs struct {
	Name            string  `pulumi:"name"`
	Enabled         *bool   `pulumi:"enabled,optional"`
	DisplayName     *string `pulumi:"displayName,optional"`
	DisplayNameHtml *string `pulumi:"displayNameHtml,optional"`
	LoginTheme      *string `pulumi:"loginTheme,optional"`
	AccountTheme    *string `pulumi:"accountTheme,optional"`
	AdminTheme      *string `pulumi:"adminTheme,optional"`
	EmailTheme      *string `pulumi:"emailTheme,optional"`
}

type RealmState struct {
	ID              string  `pulumi:"realmId"` // The ID of the realm (same as name)
	Name            string  `pulumi:"name"`
	Enabled         *bool   `pulumi:"enabled,optional"`
	DisplayName     *string `pulumi:"displayName,optional"`
	DisplayNameHtml *string `pulumi:"displayNameHtml,optional"`
	LoginTheme      *string `pulumi:"loginTheme,optional"`
	AccountTheme    *string `pulumi:"accountTheme,optional"`
	AdminTheme      *string `pulumi:"adminTheme,optional"`
	EmailTheme      *string `pulumi:"emailTheme,optional"`
}

func getClientToken(ctx context.Context) (client *gocloak.GoCloak, token *gocloak.Token) {
	config := infer.GetConfig[config.KeycloakConfig](ctx)
	client = config.Client
	token, err := client.LoginAdmin(ctx, config.AdminUsername, config.AdminPassword, config.Realm)
	return config.Client, token
}

func (*Realm) Create(ctx context.Context, req infer.CreateRequest[RealmArgs]) (infer.CreateResponse[RealmState], error) {
	if req.DryRun {
		return infer.CreateResponse[RealmState]{
			ID: req.Name,
			Output: RealmState{
				Name:            req.Inputs.Name,
				DisplayName:     req.Inputs.DisplayName,
				DisplayNameHtml: req.Inputs.DisplayNameHtml,
				Enabled:         req.Inputs.Enabled,
				LoginTheme:      req.Inputs.LoginTheme,
				AccountTheme:    req.Inputs.AccountTheme,
				AdminTheme:      req.Inputs.AdminTheme,
				EmailTheme:      req.Inputs.EmailTheme,
			},
		}, nil
	}

	client, token := getClientToken(ctx)
	realm := gocloak.RealmRepresentation{
		Realm:           &req.Inputs.Name,
		Enabled:         req.Inputs.Enabled,
		DisplayName:     req.Inputs.DisplayName,
		DisplayNameHtml: req.Inputs.DisplayNameHtml,
		LoginTheme:      req.Inputs.LoginTheme,
		AccountTheme:    req.Inputs.AccountTheme,
		AdminTheme:      req.Inputs.AdminTheme,
		EmailTheme:      req.Inputs.EmailTheme,
	}

	resp, err := client.CreateRealm(ctx, token.AccessToken, realm)
	if err != nil {
		return infer.CreateResponse[RealmState]{}, err
	}

	p.GetLogger(ctx).Info("Realm created with ID: " + resp)

	currentRealm, err := client.GetRealm(ctx, token.AccessToken, req.Inputs.Name)
	if err != nil {
		return infer.CreateResponse[RealmState]{}, err
	}

	return infer.CreateResponse[RealmState]{
		ID: resp,
		Output: RealmState{
			ID:              resp,
			Name:            currentRealm.Realm,
			DisplayName:     currentRealm.DisplayName,
			DisplayNameHtml: currentRealm.DisplayNameHtml,
			Enabled:         currentRealm.Enabled,
			LoginTheme:      currentRealm.LoginTheme,
			AccountTheme:    currentRealm.AccountTheme,
			AdminTheme:      currentRealm.AdminTheme,
			EmailTheme:      currentRealm.EmailTheme,
		},
	}, nil
}

func (r *Realm) Update(ctx context.Context, req infer.UpdateRequest[RealmArgs, RealmState]) (infer.UpdateResponse[RealmState], error) {
	if req.DryRun {
		return infer.UpdateResponse[RealmState]{
			Output: RealmState{
				Name:            req.Inputs.Name,
				DisplayName:     req.Inputs.DisplayName,
				DisplayNameHtml: req.Inputs.DisplayNameHtml,
				Enabled:         req.Inputs.Enabled,
				LoginTheme:      req.Inputs.LoginTheme,
				AccountTheme:    req.Inputs.AccountTheme,
				AdminTheme:      req.Inputs.AdminTheme,
				EmailTheme:      req.Inputs.EmailTheme,
			},
		}, nil
	}

	client, token := getClientToken(ctx)
	realm := gocloak.RealmRepresentation{
		Realm:           &req.Inputs.Name,
		Enabled:         req.Inputs.Enabled,
		DisplayName:     req.Inputs.DisplayName,
		DisplayNameHtml: req.Inputs.DisplayNameHtml,
		LoginTheme:      req.Inputs.LoginTheme,
		AccountTheme:    req.Inputs.AccountTheme,
		AdminTheme:      req.Inputs.AdminTheme,
		EmailTheme:      req.Inputs.EmailTheme,
	}

	err := client.UpdateRealm(ctx, token.AccessToken, realm)
	if err != nil {
		return infer.UpdateResponse[RealmState]{}, err
	}

	p.GetLogger(ctx).Info("Realm updated with ID: " + req.ID)

	currentRealm, err := client.GetRealm(ctx, token.AccessToken, req.Inputs.Name)
	if err != nil {
		return infer.UpdateResponse[RealmState]{}, err
	}

	return infer.UpdateResponse[RealmState]{
		Output: RealmState{
			ID:              req.ID,
			Name:            currentRealm.Realm,
			DisplayName:     currentRealm.DisplayName,
			DisplayNameHtml: currentRealm.DisplayNameHtml,
			Enabled:         currentRealm.Enabled,
			LoginTheme:      currentRealm.LoginTheme,
			AccountTheme:    currentRealm.AccountTheme,
			AdminTheme:      currentRealm.AdminTheme,
			EmailTheme:      currentRealm.EmailTheme,
		},
	}, nil
}

func (r *Realm) Read(ctx context.Context, req infer.ReadRequest[RealmArgs, RealmState]) (infer.ReadResponse[RealmArgs, RealmState], error) {

	client, token := getClientToken(ctx)
	currentRealm, err := client.GetRealm(ctx, token.AccessToken, req.ID)
	if err != nil {
		return infer.ReadResponse[RealmArgs, RealmState]{}, err
	}
	p.GetLogger(ctx).Info("Realm read with ID: " + req.ID)

	return infer.ReadResponse[RealmArgs, RealmState]{
		Output: RealmState{
			ID:              req.ID,
			Name:            currentRealm.Realm,
			DisplayName:     currentRealm.DisplayName,
			DisplayNameHtml: currentRealm.DisplayNameHtml,
			Enabled:         currentRealm.Enabled,
			LoginTheme:      currentRealm.LoginTheme,
			AccountTheme:    currentRealm.AccountTheme,
			AdminTheme:      currentRealm.AdminTheme,
			EmailTheme:      currentRealm.EmailTheme,
		},
	}, nil
}

func (r *Realm) Delete(ctx context.Context, req infer.DeleteRequest[RealmState]) (infer.DeleteResponse[RealmState], error) {

	client, token := getClientToken(ctx)
	err := client.DeleteRealm(ctx, token.AccessToken, req.ID)
	if err != nil {
		return infer.DeleteResponse[RealmState]{}, err
	}
	p.GetLogger(ctx).Info("Realm deleted with ID: " + req.ID)

	return infer.DeleteResponse[RealmState]{}, nil
}

func (r *Realm) WireDependecies(f infer.FieldSelector, args *RealmArgs, state *RealmState) {
	f.OutputField(&state.ID).DependsOnField(&args.Name)
	f.OutputField(&state.Enabled).DependsOnField(&args.Enabled)
	f.OutputField(&state.DisplayName).DependsOnField(&args.DisplayName)
	f.OutputField(&state.DisplayNameHtml).DependsOnField(&args.DisplayNameHtml)
	f.OutputField(&state.LoginTheme).DependsOnField(&args.LoginTheme)
	f.OutputField(&state.AccountTheme).DependsOnField(&args.AccountTheme)
	f.OutputField(&state.AdminTheme).DependsOnField(&args.AdminTheme)
	f.OutputField(&state.EmailTheme).DependsOnField(&args.EmailTheme)
}
