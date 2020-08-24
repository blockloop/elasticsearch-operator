package k8shandler

import (
	"context"
	"fmt"

	loggingv1 "github.com/openshift/elasticsearch-operator/pkg/apis/logging/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetElasticsearchCR(c client.Client, ns string) (*loggingv1.Elasticsearch, error) {
	esl := &loggingv1.ElasticsearchList{}
	opts := &client.ListOptions{Namespace: ns}

	if err := c.List(context.TODO(), esl, opts); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, err
		}

		return nil, fmt.Errorf("unable to get elasticsearch instance in %q: %w", ns, err)
	}

	if len(esl.Items) == 0 {
		gr := schema.GroupResource{
			Group:    loggingv1.SchemeGroupVersion.Group,
			Resource: "Elasticsearch",
		}
		return nil, apierrors.NewNotFound(gr, "elasticsearch")
	}

	return &esl.Items[0], nil
}
