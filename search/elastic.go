package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/lautistr/social-pet/enums"
	"github.com/lautistr/social-pet/models"
)

type ElasticSearchRepository struct {
	client *elastic.Client
}

func NewElastic(url string) (*ElasticSearchRepository, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, err
	}

	return &ElasticSearchRepository{
		client: client,
	}, nil
}

func (r *ElasticSearchRepository) Close() {
	//
}

func (r *ElasticSearchRepository) IndexPost(ctx context.Context, post models.Post) error {
	body, _ := json.Marshal(post)
	_, err := r.client.Index(
		enums.Elastic_PostsIndex,
		bytes.NewReader(body),
		r.client.Index.WithDocumentID(post.ID),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithRefresh("wait_for"),
	)
	return err
}

func (r *ElasticSearchRepository) SearchPost(ctx context.Context, query string) (results []models.Post, err error) {
	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":            query,
				"fields":           []string{"body"},
				"fuzzines":         3,
				"cutoff_frequency": 0.0001,
			},
		},
	}
	if err = json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(enums.Elastic_PostsIndex),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			results = nil
		}
	}()

	if res.IsError() {
		return nil, errors.New("elasticsearch error " + res.String())
	}

	var eRes map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	var posts []models.Post
	for _, hit := range eRes["hits"].(map[string]interface{})["hits"].([]interface{}) {
		post := models.Post{}
		source := hit.(map[string]interface{})["_source"]
		marshal, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(marshal, &post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func SetUpElasticSearchRepository(address string) {
	es, err := NewElastic(fmt.Sprintf("http://%s", address))
	if err != nil {
		log.Fatal(err)
	}
	SetSearchRepository(es)

	defer Close()
}
