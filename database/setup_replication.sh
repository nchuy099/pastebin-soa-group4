#!/bin/bash

# Thông tin kết nối tới MySQL Master
MASTER_HOST="mysql-primary"
MASTER_PORT="3306"
MASTER_USER="root"
MASTER_PASS="password"

# Thông tin kết nối tới MySQL Replica
REPLICA1_HOST="mysql-replica1"
REPLICA2_HOST="mysql-replica2"
REPLICA_PORT="3306"
REPLICA_USER="root"
REPLICA_PASS="password"

# Tạo user replica trên Master với mysql_native_password
echo "Creating replica user on Master using mysql_native_password..."
mysql -u"${MASTER_USER}" -p"${MASTER_PASS}" -h"${MASTER_HOST}" -P"${MASTER_PORT}" -e "
CREATE USER IF NOT EXISTS 'replica'@'%' IDENTIFIED WITH mysql_native_password BY 'replica_pass';
GRANT REPLICATION SLAVE ON *.* TO 'replica'@'%';
FLUSH PRIVILEGES;
"

# Lấy thông tin binlog từ Master
MASTER_STATUS=$(mysql -u"${MASTER_USER}" -p"${MASTER_PASS}" -h"${MASTER_HOST}" -P"${MASTER_PORT}" -e "SHOW MASTER STATUS\G")

# Trích xuất tên file binlog và vị trí log
MASTER_LOG_FILE=$(echo "${MASTER_STATUS}" | grep "File" | awk '{print $2}')
MASTER_LOG_POS=$(echo "${MASTER_STATUS}" | grep "Position" | awk '{print $2}')

echo "MASTER_LOG_FILE: ${MASTER_LOG_FILE}"
echo "MASTER_LOG_POS: ${MASTER_LOG_POS}"

# Kiểm tra nếu các biến cần thiết bị rỗng
if [ -z "${MASTER_LOG_FILE}" ] || [ -z "${MASTER_LOG_POS}" ]; then
  echo "ERROR: Could not obtain Master log file or position. Exiting..."
  exit 1
fi

# Cấu hình Replica1
echo "Configuring Replica1..."
mysql -u"${REPLICA_USER}" -p"${REPLICA_PASS}" -h"${REPLICA1_HOST}" -P"${REPLICA_PORT}" -e "
STOP SLAVE;
CHANGE MASTER TO
  MASTER_HOST='${MASTER_HOST}',
  MASTER_PORT=${MASTER_PORT},
  MASTER_USER='replica',
  MASTER_PASSWORD='replica_pass',
  MASTER_LOG_FILE='${MASTER_LOG_FILE}',
  MASTER_LOG_POS=${MASTER_LOG_POS};
START SLAVE;
"

# Cấu hình Replica2
echo "Configuring Replica2..."
mysql -u"${REPLICA_USER}" -p"${REPLICA_PASS}" -h"${REPLICA2_HOST}" -P"${REPLICA_PORT}" -e "
STOP SLAVE;
CHANGE MASTER TO
  MASTER_HOST='${MASTER_HOST}',
  MASTER_PORT=${MASTER_PORT},
  MASTER_USER='replica',
  MASTER_PASSWORD='replica_pass',
  MASTER_LOG_FILE='${MASTER_LOG_FILE}',
  MASTER_LOG_POS=${MASTER_LOG_POS};
START SLAVE;
"

# Kiểm tra trạng thái replication trên Replica1
echo "Checking replication status on Replica1..."
mysql -u"${REPLICA_USER}" -p"${REPLICA_PASS}" -h"${REPLICA1_HOST}" -P"${REPLICA_PORT}" -e "SHOW SLAVE STATUS\G" | grep -E "Slave_IO_Running:|Slave_SQL_Running:"

# Kiểm tra trạng thái replication trên Replica2
echo "Checking replication status on Replica2..."
mysql -u"${REPLICA_USER}" -p"${REPLICA_PASS}" -h"${REPLICA2_HOST}" -P"${REPLICA_PORT}" -e "SHOW SLAVE STATUS\G" | grep -E "Slave_IO_Running:|Slave_SQL_Running:"