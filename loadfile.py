import time
from locust import HttpUser, task

class QuickstartUser(HttpUser):
    @task
    def on_start(self):
        self.client.get("http://3.138.204.175/ram")