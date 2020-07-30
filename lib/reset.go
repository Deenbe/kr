package lib

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/pkg/errors"
)

type sequenceNumbers struct {
	MatchedShards   map[string]*kinesis.Record
	UnmatchedShards []string
}

type updatedSequenceNumber struct {
	OldSequenceNumber string
	NewSequenceNumber string
}

func Reset(t time.Time, config *Config) error {
	_, o := t.Zone()
	if o != 0 {
		t = t.UTC()
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{}))
	sn, err := findSequenceNumbers(t, config, sess)
	if err != nil {
		return errors.WithStack(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	if len(sn.MatchedShards) > 0 {
		fmt.Println("Matching Shards")
		fmt.Println("===============")
		fmt.Fprintln(w, "SHARD ID\tSEQUENCE NO\t")
		for sid, r := range sn.MatchedShards {
			fmt.Fprintf(w, "%v\t%v\t\n", sid, *r.SequenceNumber)
		}
		w.Flush()
	}

	if len(sn.UnmatchedShards) > 0 {
		fmt.Println("")
		fmt.Println("Shards without matching records")
		for _, sid := range sn.UnmatchedShards {
			fmt.Printf("%v\n", sid)
		}
	}

	if config.Update {
		fmt.Println("")
		fmt.Println("Consumer state")
		fmt.Println("=======================")
		u, err := updateConsumerState(sess, config, sn.MatchedShards)
		if err != nil {
			return errors.WithStack(err)
		}

		fmt.Fprintln(w, "SHARD ID\tOLD SEQUENCE NO\tNEW SEQUENCE NO\t")
		for sid, sn := range u {
			fmt.Fprintf(w, "%v\t%v\t%v\t\n", sid, sn.OldSequenceNumber, sn.NewSequenceNumber)
		}
		w.Flush()
	}

	return nil
}

func findSequenceNumbers(t time.Time, config *Config, sess *session.Session) (*sequenceNumbers, error) {
	k := kinesis.New(sess)
	r, err := k.ListShards(&kinesis.ListShardsInput{
		StreamName: &config.StreamName,
	})

	if err != nil {
		return nil, errors.WithStack(err)
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
			return nil, errors.WithStack(err)
		}

		iter := i.ShardIterator
		for {
			r, err := k.GetRecords(&kinesis.GetRecordsInput{
				ShardIterator: iter,
				Limit:         aws.Int64(1),
			})

			if err != nil {
				return nil, errors.WithStack(err)
			}

			if len(r.Records) == 1 {
				matched[*s.ShardId] = r.Records[0]
				break
			}

			if *r.MillisBehindLatest == 0 || r.NextShardIterator == nil {
				unmatched = append(unmatched, *s.ShardId)
				break
			}

			iter = r.NextShardIterator
		}
	}

	return &sequenceNumbers{MatchedShards: matched, UnmatchedShards: unmatched}, nil
}

func updateConsumerState(sess *session.Session, config *Config, newState map[string]*kinesis.Record) (map[string]*updatedSequenceNumber, error) {
	d := dynamodb.New(sess)
	m := make(map[string]*updatedSequenceNumber)

	for sid, r := range newState {
		i, err := d.GetItem(&dynamodb.GetItemInput{
			TableName: &config.ConsumerName,
			Key: map[string]*dynamodb.AttributeValue{
				"leaseKey": {S: &sid},
			},
		})

		if err != nil {
			return nil, errors.WithStack(err)
		}

		m[sid] = &updatedSequenceNumber{OldSequenceNumber: *r.SequenceNumber, NewSequenceNumber: *r.SequenceNumber}
		i.Item["checkpoint"].S = r.SequenceNumber

		_, err = d.PutItem(&dynamodb.PutItemInput{
			TableName: &config.ConsumerName,
			Item:      i.Item,
		})

		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return m, nil
}
