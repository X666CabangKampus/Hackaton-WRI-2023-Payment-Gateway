package flip

const ValidationToken = "$2y$13$e2TMXMauN6U0fzjNyJJE2ufTXr16/iUF9LKQjAaZDp4D3gybtXtUa"

type Service struct {
	idempotent Idempotent
}

func (s Service) Service() {
	s.idempotent = Idempotent{}
}

func ReqDisbursement() {

}
