package fcm

type Firebase interface {
	SendCloudMessage(model CloudMessage) (string, error)
}
