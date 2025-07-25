package ads

// create a new invoke id and make a channel for it
func (c *Client) getInvokeID() (uint32, chan Response) {
	c.mutex.Lock()
	c.invokeID++
	id := c.invokeID
	ch := make(chan Response)
	c.requests[id] = ch
	c.mutex.Unlock()
	c.logger.Debug("send: Assigned InvokeID", "invokeID", id)
	return id, ch
}

// Clean up the channel for the invokeID
func (c *Client) removeInvokeId(id uint32) {
	c.mutex.Lock()
	delete(c.requests, id)
	c.mutex.Unlock()
	c.logger.Debug("send: Cleaned up request", "invokeID", id)

}
