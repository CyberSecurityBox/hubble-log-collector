package utils

import (
	"errors"
	"github.com/cilium/hubble/pkg/env"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func SetClusterName(log *logrus.Logger) error {
	/* https://stackoverflow.com/questions/38242062/how-to-get-kubernetes-cluster-name-from-k8s-api */
	var url = "http://metadata.google.internal/computeMetadata/v1/instance/attributes/cluster-name"

	if env.ClusterName == "" {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.WithError(err).Error("SetClusterName cannot get cluster-name")
			return err
		}
		req.Header.Add("Metadata-Flavor", "Google")
		resp, err := client.Do(req)
		if err != nil {
			log.WithError(err).Error("SetClusterName cannot get cluster-name")
			return err
		}
		defer resp.Body.Close()
		clusterName, err := io.ReadAll(resp.Body)
		if err != nil {
			log.WithError(err).Error("SetClusterName cannot read response")
			return err
		}
		env.ClusterName = string(clusterName)
		log.Info("K8S_CLUSTER_NAME: ", env.ClusterName)
		if env.ClusterName == "" {
			err := errors.New("K8S_CLUSTER_NAME is the empty string")
			log.WithError(err)
			return err
		}
	}

	return nil
}
