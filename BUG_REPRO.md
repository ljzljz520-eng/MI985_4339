# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	example.com/nursery-cms/cmd/nursery-cms	[no test files]
?   	example.com/nursery-cms/collaboration	[no test files]
?   	example.com/nursery-cms/domain	[no test files]
ok  	example.com/nursery-cms/reporting	0.001s
--- FAIL: TestBusiness30Regression (0.01s)
    business30_regression_test.go:43: expected cancellation, got err=<nil> result={BatchID:985-30 Accepted:[{ID:record-985-30-one BatchID:985-30 Title:交接课程 Content:真实结果 Status:draft Version:1 Owner:teacher UpdatedAt:fixed}] Rejected:[] Cancelled:false Message:accepted=1 rejected=0}
FAIL
FAIL	example.com/nursery-cms/service	0.024s
ok  	example.com/nursery-cms/store	0.012s
ok  	example.com/nursery-cms/transport	0.007s
ok  	example.com/nursery-cms/validation	0.001s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/nursery-cms): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/nursery-cms): exit `0`
