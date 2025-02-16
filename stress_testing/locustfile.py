from locust import HttpUser, task, between
import random

class MerchStoreUser(HttpUser):
    wait_time = between(1, 5)

    token = None
    user_id = None
    prices = {
        "t-shirt": 80,
        "cup": 20,
        "book": 50,
        "pen": 10,
        "powerbank": 200,
        "hoody": 300,
        "umbrella": 200,
        "socks": 10,
        "wallet": 50,
        "pink-hoody": 500,
    }

    def on_start(self):
        """Simulates login for each user and stores token"""
        self.user_id = random.randint(1, 100000)
        username = f"user{self.user_id}"
        password = "password"
        with self.client.post(
            "/api/auth",
            json={"username": username, "password": password},
            catch_response=True,
        ) as response:
            if response.status_code == 200:
                try:
                    self.token = response.json()["token"]
                except KeyError:
                    response.failure(f"Token not found in response: {response.text}")
            else:
                response.failure(f"Login failed for {username}")

    @task(1)
    def get_info(self):
        headers = {"Authorization": f"Bearer {self.token}"}
        with self.client.get("/api/info", headers=headers, catch_response=True) as response:
            if response.status_code != 200:
                response.failure(f"Failed to get info: {response.status_code}")

    @task(1)
    def buy_item(self):
        item = random.choice(list(self.prices.keys()))
        headers = {"Authorization": f"Bearer {self.token}"}
        with self.client.get(f"/api/buy/{item}", headers=headers, catch_response=True) as response:
            if response.status_code != 200:
                response.failure(f"Failed to buy {item}: {response.status_code}")

    @task(1)
    def send_coin(self):
        recipient_id = random.randint(1, 100000)
        if recipient_id == self.user_id:
            recipient_id = (recipient_id % 100000) + 1
        recipient = f"user{recipient_id}"

        headers = {"Authorization": f"Bearer {self.token}"}
        payload = {"toUser": recipient, "amount": random.randint(1, 100)}

        with self.client.post("/api/sendCoin", json=payload, headers=headers, catch_response=True) as response:
            if response.status_code != 200:
                response.failure(f"Failed to send coins to {recipient}: {response.status_code}")
