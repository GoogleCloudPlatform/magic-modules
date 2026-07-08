package datalineage

import (
	"context"
	"encoding/json"
	"fmt"

	lineage "cloud.google.com/go/datacatalog/lineage/apiv1"
	"cloud.google.com/go/datacatalog/lineage/apiv1/lineagepb"
	"github.com/OpenLineage/openlineage/client/go/pkg/openlineage"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

func lineageClientFromConfig(ctx context.Context, config *transport_tpg.Config) (*lineage.Client, error) {
	switch {
	case config.ExternalCredentials != nil:
		// Assuming ExternalCredentials exposes a TokenSource
		return lineage.NewClient(
			ctx,
			option.WithTokenSource(oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: config.ExternalCredentials.IdentityToken,
				},
			)),
		)

	case config.AccessToken != "":
		return lineage.NewClient(
			ctx,
			option.WithTokenSource(
				oauth2.StaticTokenSource(
					&oauth2.Token{
						AccessToken: config.AccessToken,
					},
				),
			),
		)

	case config.Credentials != "":
		return lineage.NewClient(
			ctx,
			option.WithCredentialsJSON(
				[]byte(config.Credentials),
			),
		)

	default:
		// ADC
		return lineage.NewClient(ctx)
	}
}

func emitEvent(ctx context.Context, runEvent *openlineage.RunEvent, p *transport_tpg.Config) (*lineagepb.ProcessOpenLineageRunEventResponse, diag.Diagnostics) {
	eventJSON, err := json.Marshal(runEvent)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	parent := fmt.Sprintf("projects/%s/locations/%s", p.Project, p.Region)

	payload := map[string]any{}
	if err := json.Unmarshal(eventJSON, &payload); err != nil {
		return nil, diag.FromErr(err)
	}
	s, err := structpb.NewStruct(payload)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	client, err := lineageClientFromConfig(ctx, p)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	resp, err := client.ProcessOpenLineageRunEvent(ctx,
		&lineagepb.ProcessOpenLineageRunEventRequest{
			Parent:      parent,
			OpenLineage: s,
		})

	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("ProcessOpenLineageRunEvent: %w", err))
	}
	return resp, nil
}

func getLatestRunForProcess(ctx context.Context, conf *transport_tpg.Config, process string) (string, diag.Diagnostics) {
	client, err := lineageClientFromConfig(ctx, conf)
	if err != nil {
		return "", diag.FromErr(err)
	}

	_, pErr := client.GetProcess(ctx, &lineagepb.GetProcessRequest{
		Name: process,
	})

	if pErr != nil {
		// no process
		return "", diag.FromErr(pErr)
	}

	runs := client.ListRuns(ctx, &lineagepb.ListRunsRequest{
		Parent: process,
	})

	latestRun, rErr := runs.Next()
	if rErr != nil {
		// no runs for process, possible?
		return "", diag.FromErr(fmt.Errorf("Error retrieving latest run for process %s: %w", process, rErr))
	}

	return latestRun.Name, nil
}

func deleteProcess(ctx context.Context, conf *transport_tpg.Config, process string) error {
	client, err := lineageClientFromConfig(ctx, conf)
	if err != nil {
		return err
	}

	_, err = client.DeleteProcess(ctx, &lineagepb.DeleteProcessRequest{
		Name: process,
	})

	if err != nil {
		return err
	}

	return nil
}
