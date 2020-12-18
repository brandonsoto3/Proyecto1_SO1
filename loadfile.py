from locust import HttpLocust, TaskSet,task,between

#http://newtours.demoaut.com"

class UserBehaviour(TaskSet):
    @task
    def launch_Url(self):
        self.client.get("http://3.138.204.175/ram")


class User(HttpLocust):
    task_set=UserBehaviour
    wait_time = between(5, 10)