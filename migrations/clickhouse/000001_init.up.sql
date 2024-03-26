CREATE TABLE IF NOT EXISTS goods 
(
	id INTEGER, 
	project_id INTEGER,
 	name String,
  	description String,
   	priority INTEGER,
    removed BOOL,
 	event_time DATETIME,
 	INDEX idx_goods_name name TYPE ngrambf_v1(3, 512, 5, 12345) GRANULARITY 1
) 
ENGINE=MergeTree
PRIMARY KEY (id, project_id)
ORDER BY (id, project_id, event_time)