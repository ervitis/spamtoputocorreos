package bots

type (
	IBot interface {
		SendNotification(string) error
		StartServer() error
	}
)
