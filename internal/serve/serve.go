package serve

import (
	"argocd_app_sync_checker/internal/serve/config"
	"argocd_app_sync_checker/internal/serve/metrics"
	"log"
	"net/http"
)

func Start() {

	parsed_flags, err := config.ParseFlags()
	if err != nil {
		log.Fatalf("failed to parse flags due to error: %s", err)
	}

	config, err := config.ParseConfig(*parsed_flags)
	if err != nil {
		log.Fatalf("failed to read config due to error: %s", err)
	}

	metrics.RecordMetrics(&metrics.AutoSyncMetricsOptions{
		Address:         config.Argocd.Instance,
		IntervalSeconds: config.ScrapeInterval,
		Username:        config.Argocd.Username,
		Password:        config.Argocd.Password,
	})
	log.Fatal(http.ListenAndServe(config.ListenAddress, nil))

}
