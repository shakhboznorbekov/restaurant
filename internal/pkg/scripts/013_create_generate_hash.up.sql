-- multiple hash
DROP FUNCTION IF EXISTS generate_multiple_hash(links character varying[]);

CREATE OR REPLACE FUNCTION generate_multiple_hash(links character varying[])
    RETURNS character varying AS $$
DECLARE
hashH character varying default '';
    hashD character varying default '';
    h character varying;
    d character varying;
    val character varying;
    base_end_point character varying default 'http://restu.uz/media/';
    date_utc timestamp default now() + interval '15 minutes';
    index integer;
    new_map_value character varying;
    j integer default 1;
    chars1 character varying[] default '{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "0", "9", "8", "7", "6", "5", "4", "3", "2", "1", "/", ".", ":", "-", " "}';
    chars2 character varying[] default '{"H", "G", "$", "I", "S", "z", "O", "e", "7", "U", "M", "Y", "r", "l", "K", "m", "a", "v", "+", "C", "y", "6", "J", "-", "i", "4", "1", "F", "t", "B", "2", "p", "8", "h", "A", "g", "c", "X", "T", "u", "o", "k", "W", "V", "w", "N", "s", "P", "D", "b", "0", "x", "L", "Q", "@", "9", "Z", "3", "j", "n", "q", "f", "!", "E", "R", "5", "d"}';
    res character varying[] default '{}';
    hash_link character varying default '';
BEGIN
    IF links ISNULL THEN
        RETURN NULL;
END IF;
    -- Initialize hashH and hashD
SELECT
    TO_CHAR(date_utc, 'HH24:MI:SS'),
    TO_CHAR(date_utc, 'DD.MM.YYYY'),
    TO_CHAR(date_utc, 'SS')
INTO
    h,
    d,
    index;
index = index + 1;

    -- Generate hash for h [hour]
    FOREACH val IN ARRAY regexp_split_to_array(h, '') LOOP
            hashH = concat(hashH, (SELECT dest_char FROM char_map WHERE src_char=val));
END LOOP;

    -- Making new_map with chars1 and chars2 from index to 60
FOR i IN index.. 67 LOOP
            INSERT INTO new_map VALUES (chars1[j], chars2[i]);
            j = j + 1;
END LOOP;

    -- Making new_map with chars1 and chars2 from 1 to index
FOR i IN 1.. index-1 LOOP
            INSERT INTO new_map VALUES (chars1[j], chars2[i]);
            j = j + 1;
END LOOP;

    -- Generate hash for d [date]
    FOREACH val IN ARRAY regexp_split_to_array(d, '') LOOP
            new_map_value = (SELECT dest_char FROM new_map WHERE src_char=val);
            IF new_map_value ISNULL THEN
                hashD = concat(hashD, val);
ELSE
                hashD = concat(hashD, new_map_value);
END IF;
END LOOP;

    -- Make response with hashH, hashD and new_map
    FOREACH val IN ARRAY links LOOP
            hash_link = '';
            FOREACH val IN ARRAY regexp_split_to_array(val, '') LOOP
                    new_map_value = (SELECT dest_char FROM new_map WHERE src_char=val);
                    IF new_map_value ISNULL THEN
                        hash_link = concat(hash_link, val);
ELSE
                        hash_link = concat(hash_link, new_map_value);
END IF;
END LOOP;
            res = array_append(res, base_end_point || hash_link || (SELECT dest_char FROM new_map WHERE src_char=' ') || hashD || (SELECT dest_char FROM new_map WHERE src_char=' ') || hashH);
END LOOP ;
TRUNCATE TABLE new_map;
RETURN res;
END;
$$ LANGUAGE plpgsql;

-- single hash
DROP FUNCTION IF EXISTS generate_single_hash(link varchar);

CREATE OR REPLACE FUNCTION generate_single_hash(link character varying)
    RETURNS character varying AS $$
DECLARE
hashH character varying default '';
    hashD character varying default '';
    h character varying;
    d character varying;
    val character varying;
    base_end_point character varying default 'http://restu.uz/media/';
    date_utc timestamp default now() + interval '15 minutes';
    index integer;
    j integer default 1;
    new_map_value character varying;
    chars1 character varying[] default '{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "0", "9", "8", "7", "6", "5", "4", "3", "2", "1", "/", ".", ":", "-", " "}';
    chars2 character varying[] default '{"H", "G", "$", "I", "S", "z", "O", "e", "7", "U", "M", "Y", "r", "l", "K", "m", "a", "v", "+", "C", "y", "6", "J", "-", "i", "4", "1", "F", "t", "B", "2", "p", "8", "h", "A", "g", "c", "X", "T", "u", "o", "k", "W", "V", "w", "N", "s", "P", "D", "b", "0", "x", "L", "Q", "@", "9", "Z", "3", "j", "n", "q", "f", "!", "E", "R", "5", "d"}';
    res character varying default '';
    hash_link character varying default '';
BEGIN
    IF link ISNULL THEN
        RETURN NULL;
END IF;
    -- Initialize hashH and hashD
SELECT
    TO_CHAR(date_utc, 'HH24:MI:SS'),
    TO_CHAR(date_utc, 'DD.MM.YYYY'),
    TO_CHAR(date_utc, 'SS')
INTO
    h,
    d,
    index;
index = index + 1;

    -- Generate hash for h [hour]
    FOREACH val IN ARRAY regexp_split_to_array(h, '') LOOP
            hashH = concat(hashH, (SELECT dest_char FROM char_map WHERE src_char=val));
END LOOP;

    -- Making new_map with chars1 and chars2 from index to 60
FOR i IN index.. 67 LOOP
            INSERT INTO new_map VALUES (chars1[j], chars2[i]);
            j = j + 1;
END LOOP;

    -- Making new_map with chars1 and chars2 from 1 to index
FOR i IN 1.. index-1 LOOP
            INSERT INTO new_map VALUES (chars1[j], chars2[i]);
            j = j + 1;
END LOOP;

    -- Generate hash for d [date]
    FOREACH val IN ARRAY regexp_split_to_array(d, '') LOOP
            new_map_value = (SELECT dest_char FROM new_map WHERE src_char=val);
            IF new_map_value ISNULL THEN
                hashD = concat(hashD, val);
ELSE
                hashD = concat(hashD, new_map_value);
END IF;
END LOOP;

    -- Make response with hashH, hashD and new_map
    hash_link = '';
    FOREACH val IN ARRAY regexp_split_to_array(link, '') LOOP
            new_map_value = (SELECT dest_char FROM new_map WHERE src_char=val);
            IF new_map_value ISNULL THEN
                hash_link = concat(hash_link, val);
ELSE
                hash_link = concat(hash_link, new_map_value);
END IF;
END LOOP;
    res = base_end_point || hash_link || (SELECT dest_char FROM new_map WHERE src_char=' ') || hashD || (SELECT dest_char FROM new_map WHERE src_char=' ') || hashH;

TRUNCATE TABLE new_map;
RETURN res;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS new_map (
    src_char character varying,
    dest_char character varying
);