import time
from locust import HttpUser, task

class QuickstartUser(HttpUser):
    @task
    def access_model(self):
        self.client.get("http://3.138.204.175/ram")

    def on_start(self):
        self.client.get("http://3.138.204.175/procesos")