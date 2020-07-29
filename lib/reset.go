package lib

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

type sequenceNumbers struct {
	MatchedShards   map[string]*kinesis.Record
	UnmatchedShards []string
}

func Reset(t time.Time, config *Config) error {
	_, o := t.Zone()
	if o != 0 {
		t = t.UTC()
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{}))
	sn, err := findSequenceNumbers(t, config, sess)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	defer w.Flush()
	if len(sn.MatchedShards) > 0 {
		fmt.Println("Matching Shards")
		fmt.Println("===============")
		fmt.Fprintln(w, "SHARD ID\tSEQUENCE NUMBER\t")
		for sid, r := range sn.MatchedShards {
			fmt.Fprintf(w, "%v\t%v\t\n", sid, *r.SequenceNumber)
		}
	}

	if len(sn.UnmatchedShards) > 0 {
		fmt.Println("")
		fmt.Println("Shards without matching records")
		for _, sid := range sn.UnmatchedShards {
			fmt.Printf("%v\n", sid)
		}
	}
	return nil
}

func findSequenceNumbers(t time.Time, config *Config, sess *session.Session) (*sequenceNumbers, error) {
	k := kinesis.New(sess)
	r, err := k.ListShards(&kinesis.ListShardsInput{
		StreamName: &config.StreamName,
	})

	if err != nil {
		return nil, err
	}

	matched := make(map[string]*kinesis.Record)
	unmatched := make([]string, 0)
	for _, s := range r.Shards {
		i, err := k.GetShardIterator(&kinesis.GetShardIteratorInput{
			StreamName:        &config.StreamName,
			ShardId:           s.ShardId,
			ShardIteratorType: aws.String(kinesis.ShardIteratorTypeAtTimestamp),
			Timestamp:         &t,
		})

		if err != nil {
			return nil, err
		}

		iter := i.ShardIterator
		for {
			r, err := k.GetRecords(&kinesis.GetRecordsInput{
				ShardIterator: iter,
				Limit:         aws.Int64(1),
			})

			if err != nil {
				return nil, err
			}

			if len(r.Records) == 1 {
				matched[*s.ShardId] = r.Records[0]
				break
			}

			if *r.MillisBehindLatest == 0 {
				unmatched = append(unmatched, *s.ShardId)
				break
			}

			iter = r.NextShardIterator
		}
	}

	return &sequenceNumbers{MatchedShards: matched, UnmatchedShards: unmatched}, nil
}
