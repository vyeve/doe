package data

const (
	insertStmt = `
INSERT INTO
    ports (
        port_id,
        name,
        city,
        country,
        regions,
        coordinates,
        province,
        timezone,
        unlocs,
        code
    )
VALUES
    %s ON CONFLICT (port_id) DO
UPDATE
SET
    (
        port_id,
        name,
        city,
        country,
        regions,
        coordinates,
        province,
        timezone,
        unlocs,
        code
    ) = (
        EXCLUDED.port_id,
        EXCLUDED.name,
        EXCLUDED.city,
        EXCLUDED.country,
        EXCLUDED.regions,
        EXCLUDED.coordinates,
        EXCLUDED.province,
        EXCLUDED.timezone,
        EXCLUDED.unlocs,
        EXCLUDED.code
    );
`
	cleanupAliasesStm = `
DELETE FROM
    port_aliases
WHERE
    port_id IN (%s);	
`
	insertAliasesStmt = `
INSERT INTO
    port_aliases (port_id, alias)
VALUES
	%s
ON CONFLICT DO NOTHING;	
`
	getByIDStmt = `
SELECT
    pp.port_id,
    name,
    city,
    province,
    country,
    regions,
    coordinates,
    timezone,
    unlocs,
    code,
    ARRAY_AGG(alias) AS alias
FROM
    ports pp
    LEFT OUTER JOIN port_aliases pa ON pa.port_id = pp.port_id
WHERE
    pp.port_id = $1
GROUP BY
    pp.port_id;    
`
	listAllStmt = `
SELECT
    pp.port_id,
    name,
    city,
    province,
    country,
    regions,
    coordinates,
    timezone,
    unlocs,
    code,
    ARRAY_AGG(alias) AS alias
FROM
    ports pp
    LEFT OUTER JOIN port_aliases pa ON pa.port_id = pp.port_id
GROUP BY
    pp.port_id
ORDER BY
    pp.port_id
LIMIT $1;    
`
)
