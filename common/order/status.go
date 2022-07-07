package order

type Status string

const (
	Created         Status = "created"
	Cancelled       Status = "cancelled"
	AwaitingPayment Status = "awaiting-payment"
	Complete        Status = "complete"
)

func (Status) Values() (kinds []string) {
	for _, s := range []Status{Created, Cancelled, AwaitingPayment, Complete} {
		kinds = append(kinds, string(s))
	}
	return
}
