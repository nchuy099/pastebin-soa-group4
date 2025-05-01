from locust import HttpUser, task, between
import random
import datetime  
import string
import csv
import os

class PastebinUser(HttpUser):
    # Sử dụng API Gateway URL
    host = os.environ.get("API_GATEWAY_URL", "http://nginx")
    wait_time = between(1, 5)
    

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

    @task(1)
    def create_paste(self):
        """Create a new paste via POST /create-paste/api/paste"""
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
            f"{self.host}/create-paste/api/paste",
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
                            PastebinUser.paste_ids.append(paste_id)
                    except:
                        pass
                    response.success()
                    print(f"[Success] Created new paste (status={response.status_code})")
                elif response.status_code == 0:
                    response.failure("Connection failed: No response")
                    print(f"[Error] Failed to create paste (status={response.status_code})")
                else:
                    try:
                        error_msg = response.json().get("error", "No error message")
                    except Exception:
                        error_msg = response.text or "Invalid or empty response"
                    response.failure(f"Failed with status {response.status_code}: {error_msg}")
                    print(f"[Error] Failed to create paste (status={response.status_code}): {error_msg}")
            except Exception as e:
                response.failure(f"Exception: {str(e)}")
                print(f"[Error] Exception when creating paste: {str(e)}")

    @task(10)
    def get_paste_by_id(self):
        if not PastebinUser.paste_ids:
            return

        paste_id = random.choice(PastebinUser.paste_ids)
        with self.client.get(
            f"{self.host}/get-paste/api/paste/{paste_id}",
            catch_response=True,
            name="Get Paste by ID"
        ) as response:
            try:
                if response.status_code == 200:
                    response.success()
                    print(f"[Success] Retrieved paste {paste_id} (status={response.status_code})")
                elif response.status_code == 403:
                    response.success()
                    print(f"[Success] Paste {paste_id} deleted or expired (status={response.status_code})")
                elif response.status_code == 404:
                    response.success()
                    print(f"[Success] Paste {paste_id} not found (status={response.status_code})")
                elif response.status_code == 0:
                    response.failure("Connection failed: No response")
                    print(f"[Error] Failed to get paste {paste_id} (status={response.status_code})")
                else:
                    try:
                        error_msg = response.json().get("error", "No error message")
                    except Exception:
                        error_msg = response.text or "Invalid or empty response"
                    response.failure(f"Failed with status {response.status_code}: {error_msg}")
                    print(f"[Error] Failed to get paste {paste_id} (status={response.status_code}): {error_msg}")
            except Exception as e:
                response.failure(f"Exception: {str(e)}")
                print(f"[Error] Exception when getting paste {paste_id}: {str(e)}")


    # @task(3)
    # def get_monthly_stats(self):
    #     """Get monthly stats from GET /stats/api/paste/stats?month=YYYY-MM, up to the current month"""
    #     # Get current year and month
    #     today = datetime.date.today()
    #     current_year = today.year
    #     current_month = today.month
        
    #     # Randomly generate a month that is not in the future (up to the current month)
    #     year = random.randint(2020, current_year)
        
    #     if year == current_year:
    #         # If the random year is the current year, exclude months after the current month
    #         month = random.randint(1, current_month)
    #     else:
    #         # If the year is not the current year, any month is valid
    #         month = random.randint(1, 12)
        
    #     date = f"{year}-{month:02d}"
        
    #     with self.client.get(
    #         f"{self.host}/stats/api/paste/stats?month={date}",
    #         catch_response=True,
    #         name="Get Monthly Stats"
    #     ) as response:
    #         try:
    #             if response.status_code == 200:
    #                 response.success()
    #                 print(f"[Success] Retrieved monthly stats for {date} (status={response.status_code})")
    #             elif response.status_code == 0:
    #                 response.failure("Connection failed: No response")
    #                 print(f"[Error] Failed to get stats for {date} (status={response.status_code})")
    #             else:
    #                 try:
    #                     error_msg = response.json().get("error", "No error message")
    #                 except Exception:
    #                     error_msg = response.text or "Invalid or empty response"
    #                 response.failure(f"Failed with status {response.status_code}: {error_msg}")
    #                 print(f"[Error] Failed to get stats for {date} (status={response.status_code}): {error_msg}")
    #         except Exception as e:
    #             response.failure(f"Exception: {str(e)}")
    #             print(f"[Error] Exception when getting stats for {date}: {str(e)}")
