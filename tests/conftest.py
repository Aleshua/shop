import os
import uuid
import pytest
import psycopg2
from dotenv import load_dotenv
from typing import Dict 

from .api_clients.auth_client import AuthClient
from .utils.mailhog import get_latest_confirmation_code

load_dotenv(dotenv_path='.env')

def pytest_addoption(parser):
    parser.addoption(
        "--base-url",
        action="store",
        default=os.environ.get("TEST_BASE_URL", "http://localhost:80"),
        help="Base URL"
    )
    
@pytest.fixture(scope="session")
def base_url(request) -> str:
    return request.config.getoption("--base-url")

@pytest.fixture(scope="session")
def connect_auth_database():
    conn = psycopg2.connect(
        host=os.environ.get("DATABASE_HOST", "localhost"),
        port=os.environ.get("DATABASE_PORT", 5432),
        user=os.environ.get("DATABASE_USER", "user"),
        password=os.environ.get("DATABASE_PASSWORD", "password"),
        dbname=os.environ.get("DATABASE_AUTH_NAME", "postgres"),
    )
    yield conn    
    conn.close()

@pytest.fixture(autouse=True)
def cleanup(connect_auth_database):
    yield
    with connect_auth_database.cursor() as cur:
        cur.execute("TRUNCATE users, refresh_tokens, confirm_cods RESTART IDENTITY CASCADE;")
        connect_auth_database.commit()
        
@pytest.fixture(scope="session")
def auth_client(base_url):
    return AuthClient(base_url)

@pytest.fixture(scope="function")
def random_user_data():
    unique_id = str(uuid.uuid4())[:8]
    return {
        "email": f"test_{unique_id}@example.com",
        "password": "StrongPass123!"
    }

@pytest.fixture(scope="function")
def registered_user(auth_client: AuthClient, random_user_data) -> Dict[str, str]:
    resp = auth_client.register(random_user_data["email"], random_user_data["password"])
    assert resp.status_code == 200, f"не удалось зарегестрировать пользователя: {resp.text}"
    
    resp = auth_client.confirm_email(
        resp.json()["data"]["id"],
        get_latest_confirmation_code(random_user_data["email"]),
    )
    assert resp.status_code == 200, f"не удалось подтвердить пользователя: {resp.text}"
    
    return random_user_data

@pytest.fixture(scope="function")
def logged_user(auth_client: AuthClient, random_user_data) -> Dict[str, Dict[str, str]]:
    resp = auth_client.register(random_user_data["email"], random_user_data["password"])
    assert resp.status_code == 200, f"не удалось зарегестрировать пользователя: {resp.text}"
    
    resp = auth_client.confirm_email(
        resp.json()["data"]["id"],
        get_latest_confirmation_code(random_user_data["email"]),
    )
    assert resp.status_code == 200, f"не удалось подтвердить пользователя: {resp.text}"
    
    resp = auth_client.login(
        random_user_data["email"],
        random_user_data["password"],
    )
    assert resp.status_code == 200, f"не удалось войти в аккаунт пользователя: {resp.text}"
    
    return {
        "user_data": random_user_data,
        "session_data": resp.json()["data"],
    }
