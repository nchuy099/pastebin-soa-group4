import os
import random
import string
import csv
import re
import math
from bs4 import BeautifulSoup
from locust import HttpUser, task, between

class PastebinUser(HttpUser):
    wait_time = between(1, 5)
    host = os.getenv("TARGET_HOST", "http://localhost:3000")

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.paste_ids = []
        self.total_pages = None

    def on_start(self):
        pass

    def _get_total_pages_from_response(self, response):
        # Parse pagination info from HTML
        soup = BeautifulSoup(response.text, "html.parser")
        small_text = soup.find("small", string=re.compile("tổng số"))
        
        if small_text:
            match = re.search(r"tổng số (\d+) paste", small_text.text)
            if match:
                total_pastes = int(match.group(1))
                per_page = 10  # Số paste trên mỗi trang
                return math.ceil(total_pastes / per_page)
        return 1

    def _extract_error_message(self, response):
        try:
            soup = BeautifulSoup(response.text, 'html.parser')
            error_div = soup.find(id="error")
            return error_div.text.strip() if error_div else "Unknown error"
        except:
            return "Unknown error"

    @task(3)
    def view_homepage(self):
        with self.client.get("/", name="/", catch_response=True) as response:
            if response.status_code == 200 or response.status_code == 304:
                response.success()
                print(f"[Success] View homepage (status={response.status_code})")
            elif response.status_code == 0:
                response.failure("Connection failed: No response")
                print(f"[Error] Failed to view homepage (status={response.status_code})")
            else:
                response.failure(response.text)
                print(f"[Error] Failed to view homepage (status={response.status_code})")

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
                    self.paste_ids.append(paste_id)
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
        if self.paste_ids:
            paste_id = random.choice(self.paste_ids)
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
    def view_public_pastes(self):
        page = random.randint(1, self.total_pages) if self.total_pages else 1
        with self.client.get(f"/public?page={page}", name="/public?page", catch_response=True) as response:
            if response.status_code == 200:
                self.total_pages = self._get_total_pages_from_response(response)
                response.success()
                print(f"[Success] View public pastes (status={response.status_code})")
            elif response.status_code == 0:
                response.failure("Connection failed: No response")
                print(f"[Error] Failed to load public pastes (status={response.status_code})")
            else:
                response.failure(response.text)
                print(f"[Error] Failed to load public pastes (status={response.status_code})")

    @task(2)
    def view_monthly_stats(self):
        year = random.randint(2020, 2025)
        month = f"{year}-{random.randint(1, 12):02d}"
        url = f"/stats/{month}"
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