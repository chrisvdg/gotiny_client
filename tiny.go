package gotinyclient

// ListEntries lists all tiny URL entries
func (c *Client) ListEntries() ([]TinyURL, error) {

	return nil, nil
}

// CreateEntry creates a new tiny URL entry
// Returns the created TinyUrl entry
// When called without an ID, the created ID can be found there
func (c *Client) CreateEntry(id, url string) (TinyURL, error) {

	return TinyURL{}, nil
}

// UpdateEntry updates a tiny URL entry
func (c *Client) UpdateEntry(id, url string) error {

	return nil
}

// GetEntry fetches the entry information of provided ID
func (c *Client) GetEntry(id string) (TinyURL, error) {

	return TinyURL{}, nil
}

// DeleteEntry deletes an tiny URL entry
func (c *Client) DeleteEntry(id string) error {

	return nil
}
