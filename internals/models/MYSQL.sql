MYSQL-------

SELECT JSON_ARRAYAGG(JSON_OBJECT(
  'Table', TABLE_NAME,
  'Column', COLUMN_NAME,
  'Type', DATA_TYPE,
))
FROM (
  SELECT
    TABLE_NAME,
    COLUMN_NAME,
    DATA_TYPE
  FROM
    INFORMATION_SCHEMA.COLUMNS
  WHERE
    TABLE_SCHEMA = 'classicmodels'
  ORDER BY
    TABLE_NAME
) t;

SELECT 
    TABLE_SCHEMA, 
    TABLE_NAME,
    TABLE_TYPE
FROM 
    information_schema.TABLES 
WHERE 
    TABLE_SCHEMA LIKE 'classicmodels' AND TABLE_TYPE LIKE 'VIEW';

SELECT CONCAT(
    '{"schema_name":"', TABLE_SCHEMA, 
    '",', '"table_name":"', TABLE_NAME,
    '",', '"comment":"', IF(COLUMN_COMMENT = '', 'none', COLUMN_COMMENT),  
    '",', '"column_name":"', COLUMN_NAME, 
    '",', '"data_type":"', COLUMN_TYPE, 
    '"}'
  ) AS details 
FROM 
  INFORMATION_SCHEMA.COLUMNS 
WHERE 
  TABLE_SCHEMA  IN (
    'classicmodels'
  ) 
ORDER BY 
  TABLE_SCHEMA, 
  TABLE_NAME, 
  ORDINAL_POSITION;

  POSTGRES------

  1 --
  SELECT json_agg(row_to_json(t))
FROM (
  SELECT
    TABLE_NAME as "Table",
    COLUMN_NAME as "Column",
    DATA_TYPE as "Type"
  FROM
    INFORMATION_SCHEMA.COLUMNS
  WHERE
    TABLE_SCHEMA = 'public'
  ORDER BY
    TABLE_NAME
) t;

---foreignkey
SELECT
    tc.table_schema, 
    tc.constraint_name, 
    tc.table_name, 
    kcu.column_name, 
    ccu.table_schema AS foreign_table_schema,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name 
FROM information_schema.table_constraints AS tc 
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
    AND tc.table_schema = kcu.table_schema
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
WHERE tc.table_schema='classicmodels';
2 --

  SELECT 
  CONCAT(
    '{', '"schema_name": "', c.table_schema, 
    '", ', '"table_name": "', c.table_name, 
    '", ', '"comment": "', COALESCE(pg_catalog.col_description(fc.oid, c.ordinal_position), 'none'),
    '", ', '"column_name": "', c.column_name, 
    '", ', '"data_type": "', c.data_type,
    '"', '}'
  ) AS json_obj 
FROM 
  information_schema.columns c
LEFT JOIN 
  pg_catalog.pg_class fc ON fc.relname = c.table_name
LEFT JOIN 
  pg_catalog.pg_namespace ns ON ns.oid = fc.relnamespace AND ns.nspname = c.table_schema
WHERE 
  c.table_schema NOT IN ('information_schema', 'pg_catalog')
  AND c.table_schema NOT LIKE 'pg_%'
ORDER BY 
  c.table_schema, 
  c.table_name, 
  c.ordinal_position;


    SQL SERVER-----

    SELECT 
  '{"schema_name":"' + TABLE_SCHEMA + '",
"table_name":"' + TABLE_NAME + '",
"column_name":"' + COLUMN_NAME + '",
"data_type":"' + DATA_TYPE + '"
}' AS json_obj 
FROM 
  INFORMATION_SCHEMA.COLUMNS 
WHERE 
  TABLE_SCHEMA NOT IN (
    'guest', 'INFORMATION_SCHEMA', 'sys'
  ) 
ORDER BY 
  TABLE_SCHEMA, 
  TABLE_NAME, 
  ORDINAL_POSITION;

  SQLLITE 
  SELECT '{' ||
'"schema_name":"sqlite”,' ||
'"table_name":"' || m.name || '",' ||
'"column_name":"' || c.name || '",' ||
'"data_type":"' || c.type || '"' ||
'}' as json_obj
FROM sqlite_master AS m
JOIN pragma_table_info(m.name) AS c
WHERE m.type IN ('table', 'view')
ORDER BY m.name, c.cid;

SNOWFLAKE----

SELECT '{' ||
           '"schema_name":"' || TABLE_SCHEMA || '",' ||
           '"table_name":"' || TABLE_NAME || '",' ||
           '"comment":"' || COALESCE(REPLACE(NULLIF(COMMENT, ''), '''', ''''''), 'none') || '",' ||
           '"column_name":"' || COLUMN_NAME || '",' ||
           '"data_type":"' || DATA_TYPE || '"' ||
       '}' as json_obj
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA NOT IN (
  'INFORMATION_SCHEMA', 'SNOWFLAKE'
  )
ORDER BY TABLE_SCHEMA, TABLE_NAME, ORDINAL_POSITION;

ORACLE-----
SELECT 
  '{' || '"schema_name":"' || CAST(
    owner AS VARCHAR2(255)
  ) || '",' || '"table_name":"' || CAST(
    table_name AS VARCHAR2(255)
  ) || '",' || '"column_name":"' || CAST(
    column_name AS VARCHAR2(255)
  ) || '",' || '"data_type":"' || CAST(
    data_type AS VARCHAR2(255)
  ) || '"' || '}' as json_obj 
FROM 
  all_tab_cols 
WHERE 
  owner NOT IN ('SYS', 'SYSTEM') 
ORDER BY 
  owner, 
  table_name, 
  column_id;

  BIGQUERY-----

  SELECT 
  CONCAT(
    '{',
      '"schema_name":"', c.table_schema, '",',
      '"table_name":"', c.table_name, '",',
      '"comment":"',  CASE WHEN COALESCE(cf.description, '') = '' THEN 'none' ELSE cf.description END, '",',
      '"column_name":"', c.column_name, '",',
      '"data_type":"', c.data_type, '"',
    '}'
  ) AS json_string
FROM 
  `{your_dataset_name}.INFORMATION_SCHEMA.COLUMNS` AS c
LEFT JOIN 
  `{your_dataset_name}.INFORMATION_SCHEMA.COLUMN_FIELD_PATHS` AS cf
ON 
  c.table_schema = cf.table_schema 
  AND c.table_name = cf.table_name 
  AND c.column_name = cf.column_name 
WHERE NOT (
  c.table_schema LIKE 'pg_%' 
  OR c.table_schema = 'information_schema' 
  OR c.table_name LIKE 'pg_%' 
  OR c.table_name LIKE 'sql_%'
)
ORDER BY 
  c.table_schema, c.table_name, c.ordinal_position;

  MARIADB-----
  
  SELECT CONCAT(
    '{"schema_name":"', TABLE_SCHEMA, 
    '",', '"table_name":"', TABLE_NAME,
    '",', '"comment":"', IF(COLUMN_COMMENT = '', 'none', COLUMN_COMMENT),  
    '",', '"column_name":"', COLUMN_NAME, 
    '",', '"data_type":"', COLUMN_TYPE, 
    '"}'
  ) AS details 
FROM 
  INFORMATION_SCHEMA.COLUMNS 
WHERE 
  TABLE_SCHEMA NOT IN (
    'mysql', 'information_schema', 'performance_schema', 
    'sys'
  ) 
ORDER BY 
  TABLE_SCHEMA, 
  TABLE_NAME, 
  ORDINAL_POSITION;

  REDSHIFT-----

  SELECT 
  '{' +
  '"schema_name": "' + table_schema + '", ' +
  '"table_name": "' + table_name + '", '+
  '"column_name": "' + column_name + '", ' +
  '"data_type": "' + data_type + '"' +
  '}' AS json_obj
FROM information_schema.columns 
WHERE table_schema NOT IN ('information_schema', 'pg_catalog')
ORDER BY table_schema, table_name, ordinal_position;