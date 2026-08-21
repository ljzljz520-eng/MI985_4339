package domain

import "fmt"

func RecordID(batch string, external string) string {
	return fmt.Sprintf("record-%s-%s", batch, external)
}

func EventID(batch string, index int) string {
	return fmt.Sprintf("event-%s-%03d", batch, index)
}

func WorkflowID(batch string) string {
	return "workflow-" + batch
}

func AttachmentID(record, name string) string {
	return fmt.Sprintf("attachment-%s-%s", record, name)
}

func StableDigest(value string) string {
	sum := 0
	for _, ch := range value {
		sum = (sum*31 + int(ch)) % 1000003
	}
	return fmt.Sprintf("d%06d", sum)
}

func BatchKey(batch string) string {
	if batch == "" {
		return "unknown"
	}
	return batch
}
