package myaws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	log "github.com/sirupsen/logrus"
)

func LsAsg(m map[string]string) {
	/*
		cfg, err := config.LoadDefaultConfig(context.TODO())

		if err != nil {
			log.Fatal(err)
		}
	*/
	filters := makeFilters(m)

	r := getAsgs(filters)
	log.WithFields(
		log.Fields{
			"res": r,
		}).Debug("want ")
}

func leftJoin(ss ...[]string) []string {
	m := map[string]int{}
	for _, s := range ss {
		for _, v := range s {
			m[v]++
		}
	}
	res := []string{}
	for k, v := range m {
		if v == len(ss) {
			res = append(res, k)
		}
	}
	return res
}

func makeFilters(m map[string]string) [][]types.Filter {
	fltList := [][]types.Filter{}
	for k, v := range m {
		filters := []types.Filter{}

		filters = append(filters, types.Filter{
			Name:   aws.String("Value"),
			Values: []string{v},
		})

		filters = append(filters, types.Filter{
			Name:   aws.String("Key"),
			Values: []string{k},
		})

		fltList = append(fltList, filters)
	}

	return fltList
}

func getAsgs(fs [][]types.Filter) []string {

	lists := make([][]string, len(fs))
	for k, v := range fs {
		j, _ := json.Marshal(v)

		log.WithFields(
			log.Fields{
				"k": k,
				"v": string(j),
			}).Debug("paramters")

		res, err := describeTags(v)

		if err != nil {
			log.WithFields(
				log.Fields{
					"error=": err,
				}).Fatal()
		}

		log.WithFields(
			log.Fields{
				"list": res,
			}).Debug("response")

		lists[k] = res

	}

	return leftJoin(lists...)
}

func describeTags(f []types.Filter) ([]string, error) {
	// Load the Shared AWS Configuration (~/.aws/config)
	params := &autoscaling.DescribeTagsInput{
		Filters: f,
	}

	log.WithFields(
		log.Fields{
			"params": func() string {
				j, _ := json.Marshal(params)
				return string(j)
			}(),
		}).Debug("autoscaling describeTags Input")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	svc := autoscaling.NewFromConfig(cfg)

	p := autoscaling.NewDescribeTagsPaginator(svc, params, func(o *autoscaling.DescribeTagsPaginatorOptions) {
		o.Limit = 10
		o.StopOnDuplicateToken = true
	})

	resources := []string{}
	pageNum := 0
	for p.HasMorePages() && pageNum < 10 {
		resp, err := p.NextPage(context.TODO())
		if err != nil {
			log.Printf("error: %v", err)
			return nil, err

		}

		for _, v := range resp.Tags {
			resources = append(resources, aws.ToString(v.ResourceId))
		}
	}

	return resources, nil
}
