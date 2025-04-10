from locust import HttpUser, task, between
import random
import string
import csv
import os

class PastebinUser(HttpUser):
    totalItems = 741
    wait_time = between(1, 5)
    paste_ids = []

    def on_start(self):
        # Load paste IDs from CSV file
        csv_path = os.path.join(os.path.dirname(os.path.dirname(__file__)), "pastebin_id.csv")
        with open(csv_path, "r") as f:
            self.paste_ids = [line.strip() for line in f if line.strip()]

    @task(3)
    def create_paste(self):
        """Create a new paste via POST /api/paste"""
        content = ''.join(random.choices(string.ascii_letters + string.digits, k=200))
        data = {
            "content": content,
            "title": ''.join(random.choices(string.ascii_letters, k=10)),
            "language": random.choice(["text", "javascript", "python", "java", "cpp", "sql"]),
            "expires_in": random.choice(["1", "60", "1440", "10080", "43200", "525600"]),
            "visibility": random.choice(["public", "unlisted"])
        }

        with self.client.post(
            "http://localhost:8081/api/paste",
            json=data,
            catch_response=True
        ) as response:
            if response.status_code == 200 or response.status_code == 201:
                response.success()
            else:
                response.failure(f"Failed to create paste: {response.status_code}")

    @task(2)
    def get_paste_by_id(self):
        """Get a paste by ID from GET /api/paste/:id"""
        if not self.paste_ids:
            return

        paste_id = random.choice(self.paste_ids)
        with self.client.get(
            f"http://localhost:8082/api/paste/{paste_id}",
            catch_response=True
        ) as response:
            # Consider both 200 and 403 (expired) as success
            if response.status_code in [200, 403]:
                response.success()
            else:
                response.failure(f"Failed to get paste {paste_id}: {response.status_code}")

    @task(2)
    def get_pastes_by_page(self):
        """Get pastes by page from GET /api/paste?page="""
        page = totalItems / 10 + 1
        with self.client.get(
            f"http://localhost:8083/api/paste?page={page}",
            catch_response=True
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Failed to get page {page}: {response.status_code}")

    @task(1)
    def get_monthly_stats(self):
        """Get monthly stats from GET /api/paste/stats/YYYY-MM"""
        year = random.randint(2023, 2025)
        month = f"{random.randint(1, 12):02d}"
        date = f"{year}-{month}"
        
        with self.client.get(
            f"http://localhost:8084/api/paste/stats/{date}",
            catch_response=True
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Failed to get stats for {date}: {response.status_code}")