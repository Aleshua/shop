from typing import Tuple
import requests

def user_login() -> Tuple[str, str]:
    resp = requests.post(
        url="http://localhost/auth/register",
        json={
            "email": "example@example.com",
            "password": "strong_pass_123",
        },
    )
    resp = requests.post(
        url="http://localhost/auth/email/confirm",
        json={
            "user_id": resp.json()["data"]["id"],
            "code": "0000000000",
        },
    )
    resp = requests.post(
        url="http://localhost/auth/login",
        json={
            "email": "example@example.com",
            "password": "strong_pass_123",
        }
    )
    body = resp.json()
    return body["data"]["refresh_token"], body["data"]["access_token"]

def test_logout():
    refresh_token, _ = user_login()
    
    resp = requests.post(
        url="http://localhost/auth/logout",
        json={
            "refresh_token": refresh_token,
        },
    )
    
    assert resp.status_code == 200

def test_logout_with_wrong_token():    
    resp = requests.post(
        url="http://localhost/auth/logout",
        json={
            "refresh_token": "refresh_token",
        },
    )
    
    assert resp.status_code == 404
