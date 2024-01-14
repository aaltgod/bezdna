package db

const (
	queryUpsertService = `
INSERT INTO services(name, port, flag_regexp) VALUES($1, $2, $3)
ON CONFLICT ON CONSTRAINT services_port_key
	DO UPDATE SET 
				name=EXCLUDED.name, 
				flag_regexp=EXCLUDED.flag_regexp
`
	queryGetServiceByPort = `
SELECT name, port FROM services
WHERE port = $1`

	queryGetServices = `
SELECT name, port, flag_regexp FROM services
`
	queryInsertStreams = `
WITH cte_args AS (
	SELECT
		unnest($1::text[]) AS service_name,
		unnest($2::integer[]) AS service_port,
		unnest($3::text[]) AS text,
		unnest($4::text[]) AS flag_regexp,
		unnest($5::timestamp[]) AS started_at,
		unnest($6::timestamp[]) AS ended_at
)

INSERT INTO streams
		(
			service_name,
			service_port,
			text,
			flag_regexp,
			started_at,
			ended_at
		) 
		SELECT * FROM cte_args
RETURNING id
`
	queryGetStreamsByService = `
SELECT 
			id,
			service_name,
			service_port,
			text,
			flag_regexp,
			started_at,
			ended_at
 FROM streams
	WHERE service_name=$1 AND service_port=$2
			OFFSET $3 LIMIT $4	
`
	queryGetLastStreams = `
SELECT 
			id,
			service_name,
			service_port,
			text,
			flag_regexp,
			started_at,
			ended_at
FROM streams
		ORDER BY id DESC
		LIMIT $1
`
	queryGetStreams = `
SELECT 
			id,
			service_name,
			service_port,
			text,
			flag_regexp,
			started_at,
			ended_at
FROM streams
	WHERE id > $1
		LIMIT $2
`
	queryGetRegexpsByServices = `
SELECT regexp FROM flag_regexps WHERE service_name IN ($1) AND service_port IN ($2)
`
	queryUpsertFlagRegexp = `
INSERT INTO flag_regexps(regexp, service_name, service_port)
VALUES($1, $2, $3)
	ON CONFLICT ON CONSTRAINT unique_service_name_service_port 
		DO UPDATE SET regexp=EXCLUDED.regexp
`
	queryInsertFlags = `
WITH cte_args AS (
	SELECT
		unnest($1::bigint[]) AS stream_id,
		unnest($2::text[]) AS text,
		unnest($3::flag_direction[]) AS direction
)

INSERT INTO flags
		(
			stream_id,
			text,
			direction
		) 
		SELECT * FROM cte_args`
	queryGetFlagsByStreamIDs = `
SELECT 
			id,
			stream_id,
			text,
			direction
FROM flags
WHERE stream_id = ANY($1)
`
)
