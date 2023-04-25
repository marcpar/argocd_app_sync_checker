package metrics

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/session"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type AutoSyncMetricsOptions struct {
	Address         string
	IntervalSeconds uint64
	Username        string
	Password        string
}

func RecordMetrics(autoSyncOptions *AutoSyncMetricsOptions) {
	http.Handle("/metrics", promhttp.Handler())
	recordAutoSyncMetrics(autoSyncOptions)
}

func recordAutoSyncMetrics(autoSyncOptions *AutoSyncMetricsOptions) {
	go func() {

		argo_client, err := apiclient.NewClient(&apiclient.ClientOptions{
			ServerAddr: autoSyncOptions.Address,
			Insecure:   true,
		})

		if err != nil {
			log.Fatalf("failed to create argo_client due to error: %s", err)
		}

		ctx := context.Background()

		session_closer, session_client, err := argo_client.NewSessionClient()
		if err != nil {
			log.Fatalf("failed to create session client due to error: %s", err)
		}
		defer session_closer.Close()

		session_response, err := session_client.Create(ctx, &session.SessionCreateRequest{
			Username: autoSyncOptions.Username,
			Password: autoSyncOptions.Password,
		})
		if err != nil {
			log.Fatalf("failed to create session due to error: %s", err)
		}

		argo_client_with_token, err := apiclient.NewClient(&apiclient.ClientOptions{
			ServerAddr: autoSyncOptions.Address,
			AuthToken:  session_response.Token,
			Insecure:   true,
		})
		if err != nil {
			log.Fatalf("failed to create api client with token due to error: %s", err)
		}

		closer, app_client, err := argo_client_with_token.NewApplicationClient()
		if err != nil {
			log.Fatalf("failed to create application client due to error: %s", err)
		}
		defer closer.Close()

		autosync_enabled_gauge_vec := promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "argocd_autosync_enabled",
			Help: "1 if autosync is enabled on an application, else returns 0",
		}, []string{
			"application",
			"uid",
		})

		for {

			log.Print("scraping autosync status")
			apps, err := app_client.List(ctx, &application.ApplicationQuery{})
			if err != nil {
				log.Printf("failed to get application list due to error: %s", err)
				time.Sleep(time.Second * time.Duration(5))
				continue
			}

			for _, app := range apps.Items {

				gauge, err := autosync_enabled_gauge_vec.GetMetricWith(prometheus.Labels{
					"application": app.Name,
					"uid":         string(app.UID),
				})

				if err != nil {
					log.Fatalf("failed to get autosync_enabled gauge due to error: %s", err)
				}

				if app.Spec.SyncPolicy.Automated == nil {
					gauge.Set(0)
					continue
				}
				gauge.Set(1)
			}
			time.Sleep(time.Second * time.Duration(autoSyncOptions.IntervalSeconds))
		}
	}()
	log.Print("autosync scrape routine started")

}
