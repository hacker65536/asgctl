package myaws

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
)

/*
func LsAsg(){

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
*/

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

func describeTags() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	svc := autoscaling.NewFromConfig(cfg)

	params := &autoscaling.DescribeTagsInput{
		Filters: []types.Filter{
			{
				Name: aws.String("Key"),
				Values: []string{
					"Stage",
				},
			},
			{
				Name: aws.String("Value"),
				Values: []string{
					"production",
				},
			},
		},
	}

	p := autoscaling.NewDescribeTagsPaginator(svc, params, func(o *autoscaling.DescribeTagsPaginatorOptions) {
		o.Limit = 10
		o.StopOnDuplicateToken = true
	})

	pageNum := 0
	for p.HasMorePages() && pageNum < 10 {
		resp, err := p.NextPage(context.TODO())
		if err != nil {
			log.Printf("error: %v", err)
			return

		}

		for _, v := range resp.Tags {
			fmt.Println(aws.ToString(v.ResourceId))
		}
	}
}
