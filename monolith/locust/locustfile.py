import os
import random
import datetime
import string
import csv
import re
import math
from bs4 import BeautifulSoup
from locust import HttpUser, task, between

class PastebinUser(HttpUser):
    wait_time = between(1, 5)
    host = os.getenv("TARGET_HOST", "http://localhost:3000")

    # Class-level variable shared across all users
    paste_ids = []

    @classmethod
    def on_start(cls):
        """Load paste IDs from the CSV file once for all users"""
        if not cls.paste_ids:  # Only load if paste_ids are not already loaded
            csv_path = 'paste_id.csv'
            with open(csv_path, "r") as f:
                cls.paste_ids = [line.strip() for line in f if line.strip()]
            
            # Optional: Log the number of paste IDs loaded
            print(f"Loaded {len(cls.paste_ids)} paste IDs.")

    def _extract_error_message(self, response):
        try:
            soup = BeautifulSoup(response.text, 'html.parser')
            error_div = soup.find(id="error")
            return error_div.text.strip() if error_div else "Unknown error"
        except:
            return "Unknown error"

    @task(3)
    def create_paste(self):
        content = ''.join(random.choices(string.ascii_letters + string.digits, k=200))
        title = ''.join(random.choices(string.ascii_letters, k=10)) if random.random() > 0.5 else ""
        language = random.choice(["text", "javascript", "python", "java", "cpp", "sql"])
        
        # Tạo phân phối xác suất giảm dần cho thời gian hết hạn
        # Các giá trị: "" (không hết hạn), "1" (1 phút), "60" (1 giờ), "1440" (1 ngày), "10080" (1 tuần), "43200" (1 tháng)
        # Xác suất giảm dần: không hết hạn (5%), 1 phút (40%), 1 giờ (25%), 1 ngày (15%), 1 tuần (10%), 1 tháng (5%)
        expires_options = ["", "1", "60", "1440", "10080", "43200"]
        expires_weights = [5, 40, 25, 15, 10, 5]  # Tổng = 100%
        expires_in = random.choices(expires_options, weights=expires_weights, k=1)[0]
        
        visibility = random.choice(["public", "unlisted"]) if random.random() > 0.5 else "public"

        form_data = {
            "content": content,
            "title": title,
            "language": language,
            "expires_in": expires_in,
            "visibility": visibility
        }

        with self.client.post(
            "/paste", 
            data=form_data, 
            headers={"Content-Type": "application/x-www-form-urlencoded"}, 
            allow_redirects=False, 
            catch_response=True
        ) as response:
            if response.status_code == 302:
                location = response.headers.get('Location')
                if location:
                    paste_id = location.split('/')[-1]
                    PastebinUser.paste_ids.append(paste_id)
                response.success()
                print(f"[Success] Created new paste (status={response.status_code})")
            elif response.status_code == 0:
                response.failure("Connection failed: No response")
                print(f"[Error] Failed to create paste (status={response.status_code})")
            else:
                response.failure(response.text)
                print(f"[Error] Failed to create paste (status={response.status_code})")
            
    @task(3)
    def view_paste_by_id(self):
        if PastebinUser.paste_ids:
            paste_id = random.choice(PastebinUser.paste_ids)
            with self.client.get(f"/paste/{paste_id}", name="/paste/:id", catch_response=True) as response:
                if response.status_code == 200 or response.status_code == 404:
                    response.success()
                    print(f"[Success] View paste (status={response.status_code})")
                elif response.status_code == 0:
                    response.failure("Connection failed: No response")
                    print(f"[Error] Failed to view paste (status={response.status_code})")
                else:
                    response.failure(response.text)
                    print(f"[Error] Failed to view paste (status={response.status_code})")


    @task(2)
    def view_monthly_stats(self):
        """Get monthly stats from GET /stats/api/paste/stats?month=YYYY-MM, up to the current month"""
        # Get current year and month
        today = datetime.date.today()
        current_year = today.year
        current_month = today.month
        
        # Randomly generate a month that is not in the future (up to the current month)
        year = random.randint(2020, current_year)
        
        if year == current_year:
            # If the random year is the current year, exclude months after the current month
            month = random.randint(1, current_month)
        else:
            # If the year is not the current year, any month is valid
            month = random.randint(1, 12)
        
        date = f"{year}-{month:02d}"
        url = f"/stats/{date}"
        with self.client.get(url, name="/stats", catch_response=True) as response:
            if response.status_code == 200:
                response.success()
                print(f"[Success] View stats (status={response.status_code})")
            elif response.status_code == 0:
                response.failure("Connection failed: No response")
                print(f"[Error] Failed to load stats (status={response.status_code})")
            else:
                response.failure(response.text)
                print(f"[Error] Failed to load stats (status={response.status_code})")