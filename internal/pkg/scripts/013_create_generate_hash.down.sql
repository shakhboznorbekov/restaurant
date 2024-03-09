drop table if exists new_map;

-- multiple hash
DROP FUNCTION IF EXISTS generate_multiple_hash(links character varying[]);

-- single hash
DROP FUNCTION IF EXISTS generate_single_hash(link varchar);