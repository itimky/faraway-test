package book

type Book struct {
	intner intner
	quotes []string
}

func NewBook(
	intner intner,
	quotes []string,
) *Book {
	return &Book{
		intner: intner,
		quotes: quotes,
	}
}

func (b *Book) GetRandomQuote() (string, error) {
	rndIndex := b.intner.Intn(len(b.quotes))

	return b.quotes[rndIndex], nil
}
