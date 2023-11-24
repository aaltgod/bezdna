package db

const (
	queryInsertService = `
INSERT INTO services(name, port) VALUES($1, $2)
`
	queryGetServiceByPort = `
SELECT name, port FROM services
WHERE port = $1`

	queryGetServices = `
SELECT name, port FROM services
`
	queryInsertStreams = `
WITH cte_args AS (
	SELECT
		unnest($1::text[]) AS service_name,
		unnest($2::integer[]) AS service_port,
		unnest($3::text[]) AS text,
		unnest($4::timestamp[]) AS started_at,
		unnest($5::timestamp[]) AS ended_at
)

INSERT INTO streams
		(
			service_name,
			service_port,
			text,
			started_at,
			ended_at
		) 
		SELECT * FROM cte_args
`
	queryGetStreamsByService = `
SELECT 
			id,
			service_name,
			service_port,
			text,
			started_at,
			ended_at
 FROM streams
	WHERE service_name=$1 AND service_port=$2
		ORDER BY started_at DESC
			OFFSET $3 LIMIT $4	
`
)
