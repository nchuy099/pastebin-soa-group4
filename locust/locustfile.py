import os
import random
import string
from locust import HttpUser, task, between

class PastebinUser(HttpUser):
    host = None 
    wait_time = between(1, 5)
    

    HOST_CREATE = os.environ.get("HOST_CREATE", "http://localhost:8081")
    HOST_GET_ID = os.environ.get("HOST_GET_ID", "http://localhost:8082")
    HOST_GET_PAGE = os.environ.get("HOST_GET_PAGE", "http://localhost:8083")
    HOST_STATS = os.environ.get("HOST_STATS", "http://localhost:8084")

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.paste_ids = []
        self.total_pages = None  # Số lượng trang sẽ được cập nhật từ response

    def on_start(self):
        pass

    @task(3)
    def create_paste(self):
        """Create a new paste via POST /api/paste"""
        content = ''.join(random.choices(string.ascii_letters + string.digits, k=200))
        title = ''.join(random.choices(string.ascii_letters, k=10)) if random.random() > 0.5 else ""
        language = random.choice(["text", "javascript", "python", "java", "cpp", "sql"])
        
        # Tạo phân phối xác suất giảm dần cho thời gian hết hạn
        # Các giá trị: "" (không hết hạn), "1" (1 phút), "60" (1 giờ), "1440" (1 ngày), "10080" (1 tuần), "43200" (1 tháng)
        # Xác suất giảm dần: không hết hạn (5%), 1 phút (40%), 1 giờ (25%), 1 ngày (15%), 1 tuần (10%), 1 tháng (5%)
        expires_options = ["", "1", "60", "1440", "10080", "43200"]
        expires_weights = [5, 40, 25, 15, 10, 5]  # Tổng = 100%
        expires_in = random.choices(expires_options, weights=expires_weights, k=1)[0]
        
        visibility = random.choice(["PUBLIC", "UNLISTED"]) if random.random() > 0.5 else "PUBLIC"

        data = {
            "content": content,
            "title": title,
            "language": language,
            "expires_in": expires_in,
            "visibility": visibility
        }

        with self.client.post(
            f"{self.HOST_CREATE}/api/paste",
            json=data,
            catch_response=True,
            name="Create Paste"
        ) as response:
            try:
                if response.status_code == 200 or response.status_code == 201:
                    # Nếu API trả về ID của paste mới, thêm vào danh sách
                    try:
                        paste_id = response.json().get('data').get('id')
                        if paste_id:
                            self.paste_ids.append(paste_id)
                    except:
                        pass
                    response.success()
                    print(f"[Success] Created new paste (status={response.status_code})")
                elif response.status_code == 0:
                    response.failure("Connection failed: No response")
                    print(f"[Error] Failed to create paste (status={response.status_code})")
                else:
                    response.failure(f"Failed with status {response.status_code}: {response.json().get('error')}")
                    print(f"[Error] Failed to create paste (status={response.status_code}): {response.json().get('error')}")
            except Exception as e:
                response.failure(f"Exception: {str(e)}")
                print(f"[Error] Exception when creating paste: {str(e)}")

    @task(3)
    def get_paste_by_id(self):
        """Get a paste by ID from GET /api/paste/:id"""
        print(f"[Info] Current paste IDs: {self.paste_ids}")
        if not self.paste_ids:
            return

        paste_id = random.choice(self.paste_ids)
        with self.client.get(
            f"{self.HOST_GET_ID}/api/paste/{paste_id}",
            catch_response=True,
            name="Get Paste by ID"
        ) as response:
            try:
                if response.status_code == 200:
                    response.success()
                    print(f"[Success] Retrieved paste {paste_id} (status={response.status_code})")
                elif response.status_code == 403:
                    # Paste đã bị xóa hoặc đã hết hạn, vẫn coi là thành công
                    response.success()
                    print(f"[Success] Paste {paste_id} deleted or expired (status={response.status_code})")
                elif response.status_code == 404:
                    # Paste không tồn tại hoặc đã hết hạn, vẫn coi là thành công
                    response.success()
                    print(f"[Success] Paste {paste_id} not found (status={response.status_code})")
                elif response.status_code == 0:
                    response.failure("Connection failed: No response")
                    print(f"[Error] Failed to get paste {paste_id} (status={response.status_code})")
                else:
                    response.failure(f"Failed with status {response.status_code}: {response.json().get('error')}")
                    print(f"[Error] Failed to get paste {paste_id} (status={response.status_code}): {response.json().get('error')}")
            except Exception as e:
                response.failure(f"Exception: {str(e)}")
                print(f"[Error] Exception when getting paste {paste_id}: {str(e)}")

    @task(2)
    def get_public_pastes(self):
        """Get public pastes from GET /api/paste?page="""
        page = random.randint(1, self.total_pages) if self.total_pages else 1
        with self.client.get(
            f"{self.HOST_GET_PAGE}/api/paste?page={page}",
            catch_response=True,
            name="Get Public Pastes"
        ) as response:
            try:
                if response.status_code == 200:
                    # Cập nhật tổng số trang từ response
                    try:
                        data = response.json()
                        if isinstance(data, dict) and 'data' in data:
                            paste_list = data['data']
                            if isinstance(paste_list, dict) and 'pagination' in paste_list:
                                pagination = paste_list['pagination']
                                if 'totalPages' in pagination:
                                    self.total_pages = pagination['totalPages']
                    except Exception as e:
                        print(f"[Warning] Failed to update total pages: {str(e)}")
                    response.success()
                    print(f"[Success] Get public pastes page {page} (status={response.status_code})")
                elif response.status_code == 0:
                    response.failure("Connection failed: No response")
                    print(f"[Error] Failed to get public pastes page {page} (status={response.status_code})")
                else:
                    response.failure(f"Failed with status {response.status_code}: {response.json().get('error')}")
                    print(f"[Error] Failed to get public pastes page {page} (status={response.status_code}): {response.json().get('error')}")
            except Exception as e:
                response.failure(f"Exception: {str(e)}")
                print(f"[Error] Exception when getting page {page}: {str(e)}")

    @task(2)
    def get_monthly_stats(self):
        """Get monthly stats from GET /api/paste/stats?month=YYYY-MM"""
        year = random.randint(2020, 2025)
        month = f"{random.randint(1, 12):02d}"
        date = f"{year}-{month}"
        
        with self.client.get(
            f"{self.HOST_STATS}/api/paste/stats?month={date}",
            catch_response=True,
            name="Get Monthly Stats"
        ) as response:
            try:
                if response.status_code == 200:
                    response.success()
                    print(f"[Success] Retrieved monthly stats for {date} (status={response.status_code})")
                elif response.status_code == 0:
                    response.failure("Connection failed: No response")
                    print(f"[Error] Failed to get stats for {date} (status={response.status_code})")
                else:
                    response.failure(f"Failed with status {response.status_code}: {response.json().get('error')}")
                    print(f"[Error] Failed to get stats for {date} (status={response.status_code}): {response.json().get('error')}")
            except Exception as e:
                response.failure(f"Exception: {str(e)}")
                print(f"[Error] Exception when getting stats for {date}: {str(e)}")
