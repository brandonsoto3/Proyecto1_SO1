from locust import HttpUser, task, between

class MyUser(HttpUser):
    wait_time = between(5, 15)

    @task(4)
    def index(self):
        self.client.get("http://3.138.204.175/procesos")

    @task(1)
    def about(self):
        self.client.get("http://3.138.204.175/ram")
