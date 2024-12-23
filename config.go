import os
from dotenv import load_dotenv
load_dotenv()

class CDNConfig:
    def __init__(self):
        self.cache_timeout = int(os.environ.get('CDN_CACHE_TIMEOUT', 3600))
        self.server_endpoint = os.environ.get('CDN_SERVER_ENDPOINT', 'https://default-cdn.example.com')
        self.other_relevant_parameters = os.environ.get('CDN_OTHER_PARAMS', 'defaultParams')

    def display_config(self):
        print(f"CDN Cache Timeout: {self.cache_timeout}")
        print(f"CDN Server Endpoint: {self.server_endpoint}")
        print(f"Other Relevant Parameters: {self.other_relevant_parameters}")

cdn_config = CDNConfig()
cdn_config.display_config()