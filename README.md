

# squat

It's sort of like cat for SQS.

Reads from stdin and emits aws SQS messages for each line.

## usage

```
echo "hi" | squat "https://sqs.us-west-2.amazonaws.com/__XXX__/YOUR.fifo"
```

This will create an sqs event from stdin.

AWS credentials are taken from the environment or an aws credentials/config
file.

