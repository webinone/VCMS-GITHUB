package api


type APIWowzaResult struct {
	Success 	bool 		`json:"success"`
	Message  	string  	`json:"message"`
	Data 		interface{} 	`json:"data"`
}


type APIWowzaStatisticResult struct {

	ServerName			string 				`json:"serverName"`
	Uptime				int64 				`json:"uptime"`
	BytesIn				int64 				`json:"bytesIn"`
	BytesOut			int64 				`json:"bytesOut"`
	BytesInRate			int64 				`json:"bytesInRate"`
	BytesOutRate			int64 				`json:"bytesOutRate"`
	TotalConnections		int64 				`json:"totalConnections"`
	ConnectionCount			APIWowzaConnectionCount 	`json:"connectionCount"`
	ApplicationInstance		string				`json:"applicationInstance"`
	Name				string				`json:"name"`

}

type APIWowzaConnectionCount struct {
	Webm				int64 	`json:"WEBM"`
	Dvrchunks			int64 	`json:"DVRCHUNKS"`
	Rtmp				int64 	`json:"RTMP"`
	Mpegdash			int64 	`json:"MPEGDASH"`
	Cupertino			int64 	`json:"CUPERTINO"`
	Sanjose				int64 	`json:"SANJOSE"`
	Smooth				int64 	`json:"SMOOTH"`
	Rtp				int64 	`json:"RTP"`
}