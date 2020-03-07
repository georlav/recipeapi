package recipe

type QueryParams struct {
	Term        string   `json:"q"`
	Ingredients []string `json:"i"`
	Page        int64    `json:"p"`
}
