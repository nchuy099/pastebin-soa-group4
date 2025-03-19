# locust/locustfile.py
import os
from locust import HttpUser, task, between
import random
import string
from bs4 import BeautifulSoup
from db_utils import Database  # Import the Database class

class PastebinUser(HttpUser):
    wait_time = between(1, 5)  # Wait time between tasks (1-5 seconds)
    host = os.getenv("TARGET_HOST", "http://localhost:3000")  # Dynamic host
    paste_ids = []  # Store paste IDs fetched from DB

    def on_start(self):
        """Run when each user starts"""
        self.client.get("/")  # Visit homepage
        self.update_paste_ids()  # Fetch initial paste IDs from DB

    def update_paste_ids(self):
        """Update paste_ids from the database"""
        db = Database()
        self.paste_ids = db.fetch_public_paste_ids()

    @task(2)
    def view_create_form(self):
        """Simulate viewing the create paste page"""
        response = self.client.get("/")
        if response.status_code != 200:
            print(f"Failed to load create form: {response.status_code}")

    @task(3)
    def create_paste(self):
        """Simulate creating a new paste"""
        content = ''.join(random.choices(string.ascii_letters + string.digits, k=200))
        title = ''.join(random.choices(string.ascii_letters, k=10)) if random.random() > 0.5 else ""
        language = random.choice(["text", "javascript", "python", "java", "cpp", "sql"])
        expires_in = random.choice(["", "1", "60", "1440", "10080", "43200", "525600"])
        visibility = random.choice(["public", "unlisted"]) if random.random() > 0.5 else "public"

        form_data = {
            "content": content,
            "title": title,
            "language": language,
            "expires_in": expires_in,
            "visibility": visibility
        }

        response = self.client.post(
            "/paste",
            data=form_data,
            headers={"Content-Type": "application/x-www-form-urlencoded"},
            allow_redirects=False
        )

        if response.status_code == 302:  # Success redirect
            paste_url = response.headers.get("Location")
            if paste_url:
                paste_id = paste_url.split("/")[-1]
                self.client.get(paste_url, name="/paste/id")
                self.paste_ids.append(paste_id)  # Add new ID to list
        else:
            print(f"Failed to create paste: {response.status_code}")

    @task(2)
    def view_paste(self):
        """Simulate viewing a random paste"""
        if self.paste_ids:
            paste_id = random.choice(self.paste_ids)
            with self.client.get(
                f"/paste/{paste_id}",
                name="/paste/id",
                catch_response=True
            ) as response:
                if response.status_code in [200, 403]:
                    response.success()
                else:
                    response.failure(f"Failed to load paste {paste_id}: {response.status_code}")
        else:
            self.create_paste()  # Fallback if no IDs available

    @task(1)
    def view_public_pastes(self):
        """Simulate viewing the public pastes list"""
        response = self.client.get("/public")
        if response.status_code == 200:
            self.update_paste_ids()  # Refresh IDs from DB
        else:
            print(f"Failed to load public pastes: {response.status_code}")

    @task(1)
    def view_monthly_stats(self):
        """Simulate viewing monthly stats"""
        if random.random() > 0.5:
            year = random.randint(2020, 2025)
            month = f"{year}-{random.randint(1, 12):02d}"
            url = f"/stats/{month}"
        else:
            url = "/stats"
        
        response = self.client.get(url, name="/stats")
        if response.status_code != 200:
            print(f"Failed to load stats for {url}: {response.status_code}")