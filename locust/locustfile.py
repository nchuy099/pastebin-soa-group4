from locust import HttpUser, task, between
import random
import string
from bs4 import BeautifulSoup

class PastebinUser(HttpUser):
    wait_time = between(1, 5)  # Thời gian chờ giữa các task (1-5 giây)
    host = "http://localhost:8080"  # Địa chỉ ứng dụng pastebin
    paste_ids = []  # Lưu trữ danh sách ID paste // nen clear volume cua pastebin truoc khi test

    def on_start(self):
        """Chạy khi mỗi user bắt đầu"""
        self.client.get("/")
        self.update_public_paste_ids()

    def update_public_paste_ids(self):
        """Cập nhật danh sách paste IDs từ route /public"""
        response = self.client.get("/public")
        if response.status_code == 200:
            soup = BeautifulSoup(response.text, "html.parser")
            links = soup.select("a[href^='/paste/']")
            self.paste_ids = [link["href"].split("/")[-1] for link in links]
            if not self.paste_ids:
                print("No public pastes found yet")
        else:
            print(f"Failed to fetch /public: {response.status_code}")

    @task(2)
    def view_create_form(self):
        """Mô phỏng xem trang tạo paste"""
        response = self.client.get("/")
        if response.status_code != 200:
            print(f"Failed to load create form: {response.status_code}")

    @task(3)
    def create_paste(self):
        """Mô phỏng tạo paste mới"""
        content = ''.join(random.choices(string.ascii_letters + string.digits, k=200))
        title = ''.join(random.choices(string.ascii_letters, k=10)) if random.random() > 0.5 else ""
        language = random.choice(["text", "javascript", "python", "java", "cpp", "sql"])  # Khớp với JSP
        expires_in = random.choice(["", "1", "60", "1440", "10080", "43200", "525600"])  # Khớp với JSP
        visibility = random.choice(["public", "unlisted"]) if random.random() > 0.5 else "public"  # Khớp với JSP

        form_data = {
            "content": content,
            "title": title,
            "language": language,
            "expires_in": expires_in,
            "visibility": visibility
        }

        print(f"Sending paste - Title: {title}, Content: {content}, Language: {language}, Expires_in: {expires_in}, Visibility: {visibility}")

        response = self.client.post(
            "/paste",
            data=form_data,
            headers={"Content-Type": "application/x-www-form-urlencoded"},
            allow_redirects=False  # Tắt tự động redirect
        )

        if response.status_code == 302:  # Redirect khi thành công
            paste_url = response.headers.get("Location")
            if paste_url:
                paste_id = paste_url.split("/")[-1]
                self.client.get(
                    paste_url,
                    name="/paste/id"
                )
                self.paste_ids.append(paste_id)
        elif response.status_code == 200:  # Trường hợp lỗi mềm (như content rỗng)
            soup = BeautifulSoup(response.text, "html.parser")
            error = soup.select_one(".alert-danger")
            if error:
                print(f"Create paste failed with error: {error.text.strip()}")
            else:
                print("Create paste returned 200 but no redirect - unexpected behavior")
        else:
            print(f"Failed to create paste: {response.status_code}")

    @task(2)
    def view_paste(self):
        """Mô phỏng xem một paste ngẫu nhiên"""
        if self.paste_ids:
            paste_id = random.choice(self.paste_ids)
            response = self.client.get(
                f"/paste/{paste_id}",
                name="/paste/id"
            )
            if response.status_code != 200:
                print(f"Failed to load paste {paste_id}: {response.status_code}")
        else:
            self.create_paste()

    @task(1)
    def view_public_pastes(self):
        """Mô phỏng xem danh sách paste công khai"""
        response = self.client.get("/public")
        if response.status_code == 200:
            self.update_public_paste_ids()
        else:
            print(f"Failed to load public pastes: {response.status_code}")

    @task(1)
    def view_monthly_stats(self):
        """Mô phỏng xem thống kê theo tháng"""
        if random.random() > 0.5:
            year = random.randint(2020, 2025)
            month = f"{year}-{random.randint(1, 12):02d}"
            url = f"/stats/{month}"
        else:
            url = "/stats"
        
        response = self.client.get(url, name="/stats")
        if response.status_code != 200:
            print(f"Failed to load stats for {url}: {response.status_code}")