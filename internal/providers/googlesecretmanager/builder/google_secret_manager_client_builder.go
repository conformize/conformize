package builder

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/conformize/conformize/common/diagnostics"
	"github.com/conformize/conformize/serialization"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type GoogleSecretManagerBuilder struct {
	client      *secretmanager.Client
	credentials *google.Credentials
	diags       *diagnostics.Diagnostics
}

func New() *GoogleSecretManagerBuilder {
	return &GoogleSecretManagerBuilder{
		client:      nil,
		credentials: nil,
		diags:       diagnostics.NewDiagnostics(),
	}
}

func (googleSecretMngrBldr *GoogleSecretManagerBuilder) WithADC() *GoogleSecretManagerBuilder {
	googleSecretMngrBldr.credentials = &google.Credentials{}
	return googleSecretMngrBldr
}

func (googleSecretMngrBldr *GoogleSecretManagerBuilder) WithCredentialsFile(credentialsFile, scope string) *GoogleSecretManagerBuilder {
	fileSrc := serialization.FileSource{FilePath: credentialsFile}
	data, err := fileSrc.Read()
	if err != nil {
		googleSecretMngrBldr.diags.Append(diagnostics.Builder().Error().Summary(err.Error()).Build())
		return googleSecretMngrBldr
	}

	creds, err := google.CredentialsFromJSON(context.Background(), data, scope)
	if err != nil {
		googleSecretMngrBldr.diags.Append(diagnostics.Builder().Error().Summary(err.Error()).Build())
	}
	googleSecretMngrBldr.credentials = creds
	return googleSecretMngrBldr
}

func (googleSecretMngrBldr *GoogleSecretManagerBuilder) WithToken(token string) *GoogleSecretManagerBuilder {
	googleSecretMngrBldr.credentials = &google.Credentials{
		TokenSource: oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: token,
			},
		),
	}
	return googleSecretMngrBldr
}

func (googleSecretMngrBldr *GoogleSecretManagerBuilder) Build() (*secretmanager.Client, *diagnostics.Diagnostics) {
	if googleSecretMngrBldr.credentials == nil {
		googleSecretMngrBldr.diags.Append(diagnostics.Builder().
			Error().
			Summary("authentication method not specified or supported").
			Build(),
		)
		return nil, googleSecretMngrBldr.diags
	}

	var err error
	googleSecretMngrBldr.client, err = secretmanager.NewClient(
		context.Background(),
		option.WithCredentials(googleSecretMngrBldr.credentials),
	)

	if err != nil {
		googleSecretMngrBldr.diags.Append(diagnostics.Builder().
			Error().
			Summary(err.Error()).
			Build(),
		)
	}

	return googleSecretMngrBldr.client, googleSecretMngrBldr.diags
}
