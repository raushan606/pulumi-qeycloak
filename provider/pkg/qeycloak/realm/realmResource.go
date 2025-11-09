package realm

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	cfg "github.com/raushan606/pulumi-qeycloak/provider/pkg/qeycloak/config"
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

// getClientToken logs in using the configured admin credentials and returns the client & JWT.
func getClientToken(ctx context.Context) (*gocloak.GoCloak, *gocloak.JWT, error) {
	c := infer.GetConfig[cfg.KeycloakConfig](ctx)
	if c.Client == nil {
		return nil, nil, fmt.Errorf("keycloak client not configured")
	}
	realm := "master"
	if c.Realm != nil && *c.Realm != "" {
		realm = *c.Realm
	}
	jwt, err := c.Client.LoginAdmin(ctx, c.Username, c.Password, realm) // returns *gocloak.JWT
	if err != nil {
		return nil, nil, fmt.Errorf("admin login failed: %w", err)
	}
	return c.Client, jwt, nil
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

	client, token, err := getClientToken(ctx)
	if err != nil {
		return infer.CreateResponse[RealmState]{}, err
	}
	rr := gocloak.RealmRepresentation{
		Realm:           &req.Inputs.Name,
		Enabled:         req.Inputs.Enabled,
		DisplayName:     req.Inputs.DisplayName,
		DisplayNameHTML: req.Inputs.DisplayNameHtml, // map Html -> HTML
		LoginTheme:      req.Inputs.LoginTheme,
		AccountTheme:    req.Inputs.AccountTheme,
		AdminTheme:      req.Inputs.AdminTheme,
		EmailTheme:      req.Inputs.EmailTheme,
	}

	resp, err := client.CreateRealm(ctx, token.AccessToken, rr)
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
			Name:            deref(currentRealm.Realm),
			DisplayName:     currentRealm.DisplayName,
			DisplayNameHtml: currentRealm.DisplayNameHTML,
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

	client, token, err := getClientToken(ctx)
	if err != nil {
		return infer.UpdateResponse[RealmState]{}, err
	}
	rr := gocloak.RealmRepresentation{
		Realm:           &req.Inputs.Name,
		Enabled:         req.Inputs.Enabled,
		DisplayName:     req.Inputs.DisplayName,
		DisplayNameHTML: req.Inputs.DisplayNameHtml,
		LoginTheme:      req.Inputs.LoginTheme,
		AccountTheme:    req.Inputs.AccountTheme,
		AdminTheme:      req.Inputs.AdminTheme,
		EmailTheme:      req.Inputs.EmailTheme,
	}

	err = client.UpdateRealm(ctx, token.AccessToken, rr)
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
			Name:            deref(currentRealm.Realm),
			DisplayName:     currentRealm.DisplayName,
			DisplayNameHtml: currentRealm.DisplayNameHTML,
			Enabled:         currentRealm.Enabled,
			LoginTheme:      currentRealm.LoginTheme,
			AccountTheme:    currentRealm.AccountTheme,
			AdminTheme:      currentRealm.AdminTheme,
			EmailTheme:      currentRealm.EmailTheme,
		},
	}, nil
}

func (r *Realm) Read(ctx context.Context, req infer.ReadRequest[RealmArgs, RealmState]) (infer.ReadResponse[RealmArgs, RealmState], error) {

	client, token, err := getClientToken(ctx)
	if err != nil {
		return infer.ReadResponse[RealmArgs, RealmState]{}, err
	}
	currentRealm, err := readRealmState(ctx, client, token.AccessToken, req.ID)
	if err != nil {
		return infer.ReadResponse[RealmArgs, RealmState]{}, err
	}
	p.GetLogger(ctx).Info("Realm read with ID: " + req.ID)

	return infer.ReadResponse[RealmArgs, RealmState]{
		ID:     currentRealm.ID,
		Inputs: RealmArgs{
			Name:            req.Inputs.Name,
			DisplayName:     req.Inputs.DisplayName,
			DisplayNameHtml: req.Inputs.DisplayNameHtml,
			Enabled:         req.Inputs.Enabled,
			LoginTheme:      req.Inputs.LoginTheme,
			AccountTheme:    req.Inputs.AccountTheme,
			AdminTheme:      req.Inputs.AdminTheme,
			EmailTheme:      req.Inputs.EmailTheme,
		},	
		State: currentRealm,
	}, nil
}

func (r *Realm) Delete(ctx context.Context, req infer.DeleteRequest[RealmState]) (infer.DeleteResponse, error) {
	client, token, err := getClientToken(ctx)
	if err != nil {
		return infer.DeleteResponse{}, err
	}
	if err = client.DeleteRealm(ctx, token.AccessToken, req.ID); err != nil {
		return infer.DeleteResponse{}, err
	}
	p.GetLogger(ctx).Info("Realm deleted with ID: " + req.ID)
	return infer.DeleteResponse{}, nil
}

// WireDependencies informs the provider of input/output relationships.
func (r *Realm) WireDependencies(f infer.FieldSelector, args *RealmArgs, state *RealmState) {
	f.OutputField(&state.Name).DependsOn(f.InputField(&args.Name))
	f.OutputField(&state.DisplayName).DependsOn(f.InputField(&args.DisplayName))
	f.OutputField(&state.LoginTheme).DependsOn(f.InputField(&args.LoginTheme))
	f.OutputField(&state.AccountTheme).DependsOn(f.InputField(&args.AccountTheme))
	f.OutputField(&state.AdminTheme).DependsOn(f.InputField(&args.AdminTheme))
	f.OutputField(&state.EmailTheme).DependsOn(f.InputField(&args.EmailTheme))
}

// deref safely dereferences a *string returning empty string if nil.
func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func readRealmState(ctx context.Context, client *gocloak.GoCloak, token, realmName string) (RealmState, error) {
	realm, err := client.GetRealm(ctx, token, realmName)
	if err != nil {
		return RealmState{}, fmt.Errorf("failed to get realm: %w", err)
	}

	state := RealmState{
		ID:   *realm.Realm,
		Name: *realm.Realm,
	}

	if realm.Enabled != nil {
		state.Enabled = realm.Enabled
	}

	if realm.DisplayName != nil {
		state.DisplayName = realm.DisplayName
	}

	if realm.DisplayNameHTML != nil {
		state.DisplayNameHtml = realm.DisplayNameHTML
	}

	if realm.LoginTheme != nil {
		state.LoginTheme = realm.LoginTheme
	}

	if realm.AccountTheme != nil {
		state.AccountTheme = realm.AccountTheme
	}

	if realm.AdminTheme != nil {
		state.AdminTheme = realm.AdminTheme
	}

	if realm.EmailTheme != nil {
		state.EmailTheme = realm.EmailTheme
	}

	return state, nil
}
