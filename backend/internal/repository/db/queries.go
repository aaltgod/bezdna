package db

const (
	queryInsertService = `
INSERT INTO services VALUES($1, $2)
`
	queryGetServices = `
SELECT name, port FROM services
`
	queryInsertStream = `
INSERT INTO streams VALUES($1, $2, $3, $4, $5)	
`
	queryGetStreamsByService = `
SELECT ack, timestamp, payload FROM streams
	WHERE streams.service_name=$1 AND streams.service_port=$2
		ORDER BY timestamp DESC
			OFFSET $3
			
`
)
