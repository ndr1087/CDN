import os
import time
import pickle
import hashlib
from functools import wraps
from dotenv import load_dotenv

load_dotenv()

CACHE_DIR = os.getenv("CACHE_DIR", ".cache")
CACHE_EXPIRE_SECONDS = int(os.getenv("CACHE_EXPIRE_SECONDS", 3600))

if not os.path.exists(CACHE_DIR):
    os.makedirs(CACHE_DIR, exist_ok=True)

def cache_key_function(args, kwargs):
    key = "_".join(map(str, args)) + "_" + "_".join([f"{k}={v}" for k, v in kwargs.items()])
    return hashlib.md5(key.encode()).hexdigest()

def is_cache_expired(file_path, expiration_seconds):
    current_time = time.time()
    mod_time = os.path.getmtime(file_path)
    return (current_time - mod_time) > expiration_seconds

def cache(expiration_seconds=CACHE_EXPIRE_SECONDS):
    def decorator_cache(func):
        @wraps(func)
        def wrapper_cache(*args, **kwargs):
            cache_key = cache_key_function(args, kwargs)
            cache_file = os.path.join(CACHE_DIR, f"{cache_key}.cache")
            try:
                if not is_cache_expired(cache_file, expiration_seconds):
                    with open(cache_file, "rb") as cf:
                        return pickle.load(cf)
            except FileNotFoundError:
                pass
            
            result = func(*args, **kwargs)
            try:
                with open(cache_file, "wb") as cf:
                    pickle.dump(result, cf, protocol=pickle.HIGHEST_PROTOCOL)
            except:
                pass
            
            return result
        return wrapper_cache
    return decorator_cache

@cache(expiration_seconds=120)
def expensive_function(param1, param2):
    time.sleep(2)
    return f"Result of {param1} and {param2}"