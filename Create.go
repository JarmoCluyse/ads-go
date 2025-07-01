package adsgo

func Create() (client *AdsClient, err error) {

	return &AdsClient{
		Connection: AdsClientConnection{},
		Settings:   AdsClientSettings{},
	}, nil
}
