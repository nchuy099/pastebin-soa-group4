-- 1. Connect to ProxySQL admin interface
-- mysql -h 127.0.0.1 -P 6032 -u admin -padmin

-- 2. Configure monitoring user with proper credentials
UPDATE global_variables SET variable_value='monitor' WHERE variable_name='mysql-monitor_username';
UPDATE global_variables SET variable_value='StrongSecurePassword123' WHERE variable_name='mysql-monitor_password';

-- 3. Set monitoring intervals (in milliseconds)
UPDATE global_variables SET variable_value='2000' WHERE variable_name='mysql-monitor_connect_interval';
UPDATE global_variables SET variable_value='2000' WHERE variable_name='mysql-monitor_ping_interval';
UPDATE global_variables SET variable_value='2000' WHERE variable_name='mysql-monitor_read_only_interval';
UPDATE global_variables SET variable_value='5000' WHERE variable_name='mysql-monitor_replication_lag_interval';

-- 4. Add application user to mysql_users table
INSERT INTO mysql_users (username, password, default_hostgroup, active) 
VALUES ('proxysql_user', 'proxysql_password', 1, 1);

-- 5. Add backend MySQL servers (Primary and Replicas)
INSERT INTO mysql_servers (hostgroup_id, hostname, port, max_connections) 
VALUES (0, '10.148.0.2', 3306, 100);
INSERT INTO mysql_servers (hostgroup_id, hostname, port, max_connections) 
VALUES (1, '10.148.0.3', 3306, 100);
INSERT INTO mysql_servers (hostgroup_id, hostname, port, max_connections) 
VALUES (0, '10.148.0.4', 3306, 100);

-- 6. Add query rules
-- Rule 1: All SELECT queries go to replica (hostgroup 0)
INSERT INTO mysql_query_rules (rule_id, active, match_pattern, destination_hostgroup, apply) 
VALUES (1, 1, '^SELECT.*', 0, 1);

-- Rule 2: All INSERT, UPDATE, DELETE queries go to primary (hostgroup 1)
INSERT INTO mysql_query_rules (rule_id, active, match_pattern, destination_hostgroup, apply) 
VALUES (2, 1, '^(INSERT|UPDATE|DELETE).*', 1, 1);

-- 7. Additional ProxySQL settings for performance tuning
UPDATE global_variables SET variable_value='5000' WHERE variable_name='mysql-query_timeout';
UPDATE global_variables SET variable_value='100' WHERE variable_name='mysql-default_query_delay';

-- 8. Save configurations to disk
SAVE MYSQL USERS TO DISK;
SAVE MYSQL SERVERS TO DISK;
SAVE MYSQL QUERY RULES TO DISK;
SAVE MYSQL VARIABLES TO DISK;

-- 9. Load configurations to runtime
LOAD MYSQL USERS TO RUNTIME;
LOAD MYSQL SERVERS TO RUNTIME;
LOAD MYSQL QUERY RULES TO RUNTIME;
LOAD MYSQL VARIABLES TO RUNTIME;