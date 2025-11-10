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

type Args struct {
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
	client := gocloak.NewClient(c.Url)
	realm := "master"
	if c.Realm != nil && *c.Realm != "" {
		realm = *c.Realm
	}
	jwt, err := client.LoginAdmin(ctx, c.Username, c.Password, realm)
	if err != nil {
		return nil, nil, fmt.Errorf("admin login failed: %w", err)
	}
	return client, jwt, nil
}

func (*Realm) Create(ctx context.Context, req infer.CreateRequest[Args]) (infer.CreateResponse[RealmState], error) {
	if req.DryRun {
		return infer.CreateResponse[RealmState]{
			ID: req.Inputs.Name,
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
		DisplayNameHTML: req.Inputs.DisplayNameHtml,
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

func (r *Realm) Update(ctx context.Context, req infer.UpdateRequest[Args, RealmState]) (infer.UpdateResponse[RealmState], error) {
	p.GetLogger(ctx).Info("Updating Realm with ID: " + req.ID)
	if req.DryRun {
		return infer.UpdateResponse[RealmState]{
			Output: RealmState{
				ID:              req.Inputs.Name,
				Name:            req.Inputs.Name,
				Enabled:         req.Inputs.Enabled,
				DisplayName:     req.Inputs.DisplayName,
				DisplayNameHtml: req.Inputs.DisplayNameHtml,
				LoginTheme:      req.Inputs.LoginTheme,
				AccountTheme:    req.Inputs.AccountTheme,
				AdminTheme:      req.Inputs.AdminTheme,
				EmailTheme:      req.Inputs.EmailTheme,
			},
		}, nil
	}

	client, token, err := getClientToken(ctx)
	if err != nil {
		return infer.UpdateResponse[RealmState]{}, fmt.Errorf("failed to get client token: %w", err)
	}

	err = updateManagedFields(ctx, req.Inputs)
	if err != nil {
		return infer.UpdateResponse[RealmState]{}, fmt.Errorf("failed to update managed fields: %w", err)
	}

	state, err := readRealmState(ctx, client, token.AccessToken, req.Inputs.Name)
	if err != nil {
		return infer.UpdateResponse[RealmState]{}, fmt.Errorf("failed to read realm state: %w", err)
	}

	return infer.UpdateResponse[RealmState]{
		Output: state,
	}, nil
}

func (r *Realm) Read(ctx context.Context, req infer.ReadRequest[Args, RealmState]) (infer.ReadResponse[Args, RealmState], error) {

	client, token, err := getClientToken(ctx)
	if err != nil {
		return infer.ReadResponse[Args, RealmState]{}, err
	}
	currentRealm, err := readRealmState(ctx, client, token.AccessToken, req.ID)
	if err != nil {
		return infer.ReadResponse[Args, RealmState]{}, err
	}
	p.GetLogger(ctx).Info("Realm read with ID: " + req.ID)

	return infer.ReadResponse[Args, RealmState]{
		ID: currentRealm.ID,
		Inputs: Args{
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

func (r *Realm) WireDependencies(f infer.FieldSelector, args *Args, state *RealmState) {
	f.OutputField(&state.Name).DependsOn(f.InputField(&args.Name))
	f.OutputField(&state.DisplayName).DependsOn(f.InputField(&args.DisplayName))
	f.OutputField(&state.LoginTheme).DependsOn(f.InputField(&args.LoginTheme))
	f.OutputField(&state.AccountTheme).DependsOn(f.InputField(&args.AccountTheme))
	f.OutputField(&state.AdminTheme).DependsOn(f.InputField(&args.AdminTheme))
	f.OutputField(&state.EmailTheme).DependsOn(f.InputField(&args.EmailTheme))
}

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

func updateManagedFields(ctx context.Context, args Args) error {
	client, token, err := getClientToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get client token: %w", err)
	}
	currentRealm, err := client.GetRealm(ctx, token.AccessToken, args.Name)
	if err != nil {
		return fmt.Errorf("failed to get current realm: %w", err)
	}

	// Track if any managed field has changed
	hasChanges := false

	updateRealm := *currentRealm

	if args.Enabled != nil && !ptrBoolEqual(currentRealm.Enabled, args.Enabled) {
		updateRealm.Enabled = args.Enabled
		hasChanges = true
	}

	if args.DisplayName != nil && !ptrStringEqual(currentRealm.DisplayName, args.DisplayName) {
		updateRealm.DisplayName = args.DisplayName
		hasChanges = true
	}

	if args.DisplayNameHtml != nil && !ptrStringEqual(currentRealm.DisplayNameHTML, args.DisplayNameHtml) {
		updateRealm.DisplayNameHTML = args.DisplayNameHtml
		hasChanges = true
	}

	if args.LoginTheme != nil && !ptrStringEqual(currentRealm.LoginTheme, args.LoginTheme) {
		updateRealm.LoginTheme = args.LoginTheme
		hasChanges = true
	}

	if args.AccountTheme != nil && !ptrStringEqual(currentRealm.AccountTheme, args.AccountTheme) {
		updateRealm.AccountTheme = args.AccountTheme
		hasChanges = true
	}

	if args.AdminTheme != nil && !ptrStringEqual(currentRealm.AdminTheme, args.AdminTheme) {
		updateRealm.AdminTheme = args.AdminTheme
		hasChanges = true
	}

	if args.EmailTheme != nil && !ptrStringEqual(currentRealm.EmailTheme, args.EmailTheme) {
		updateRealm.EmailTheme = args.EmailTheme
		hasChanges = true
	}

	if !hasChanges {
		return nil
	}

	err = client.UpdateRealm(ctx, token.AccessToken, updateRealm)
	if err != nil {
		return fmt.Errorf("failed to update realm: %w", err)
	}

	return nil
}

func ptrBoolEqual(a, b *bool) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func ptrStringEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
