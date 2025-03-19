# locust/db_utils.py
import mysql.connector
from mysql.connector import Error
import os

class Database:
    def __init__(self):
        """Initialize database configuration from environment variables"""
        self.db_config = {
            "host": os.getenv("DB_HOST", "localhost"),
            "port": int(os.getenv("DB_PORT", "3307")),
            "user": os.getenv("DB_USER", "root"),
            "password": os.getenv("DB_PASSWORD", "your_root_password"),
            "database": os.getenv("DB_NAME", "your_db_name"),
        }

    def connect(self):
        """Establish a connection to the MySQL database"""
        try:
            connection = mysql.connector.connect(**self.db_config)
            if connection.is_connected():
                return connection
        except Error as e:
            print(f"Error connecting to MySQL: {e}")
        return None

    def fetch_public_paste_ids(self):
        """Fetch public paste IDs from the database"""
        connection = self.connect()
        paste_ids = []
        if connection:
            try:
                cursor = connection.cursor()
                query = """
                    SELECT id FROM pastes 
                    WHERE visibility = 'public' 
                    AND (expires_at IS NULL OR expires_at > NOW())
                """
                cursor.execute(query)
                rows = cursor.fetchall()
                paste_ids = [row[0] for row in rows]
                if not paste_ids:
                    print("No public pastes found in database")
                else:
                    print(f"Loaded {len(paste_ids)} paste IDs from database")
            except Error as e:
                print(f"Error fetching paste IDs: {e}")
            finally:
                cursor.close()
                connection.close()
        return paste_ids

# Example usage (for testing outside Locust)
if __name__ == "__main__":
    db = Database()
    ids = db.fetch_public_paste_ids()
    print(f"Test fetch: {ids}")