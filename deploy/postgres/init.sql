-- definition of 'ports' table
CREATE TABLE IF NOT EXISTS ports (
    port_id VARCHAR(16),
    name TEXT,
    city TEXT,
    province TEXT,
    country TEXT,
    regions TEXT [],
    coordinates DECIMAL [],
    timezone TEXT,
    unlocs TEXT [],
    code VARCHAR(16),
    PRIMARY KEY (port_id)
);

-- storage for port's aliases
CREATE TABLE IF NOT EXISTS port_aliases (
    port_id VARCHAR(16),
    alias TEXT,
    PRIMARY KEY (port_id, alias)
);