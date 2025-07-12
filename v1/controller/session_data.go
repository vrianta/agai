package controller

/*
 * Store Data in the Session
 */
func (_c *Context) StoreData(index string, _d any) {
	_c.session.Store[index] = _d
}

/*
 * Get Data From Session Store
 */
func (_c *Context) GetStoredData(index string) any {
	return _c.session.Store[index]
}
