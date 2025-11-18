import time
import requests

def test_register():
    resp = requests.post(
        url="http://localhost/auth/register",
        json={
            "email": "example@example.com",
            "password": "strong_pass_123",
        },
    )
    
    assert resp.status_code == 200
    
    body = resp.json()
    
    assert isinstance(body, dict)
    assert isinstance(body["data"], dict)
    assert isinstance(body["data"]["id"], int)

def test_register_repeat():
    requests.post(
        url="http://localhost/auth/register",
        json={
            "email": "example@example.com",
            "password": "strong_pass_123",
        },
    )
    
    resp = requests.post(
        url="http://localhost/auth/register",
        json={
            "email": "example@example.com",
            "password": "strong_pass_123",
        },
    )
    assert resp.status_code == 200

def test_register_confirm_email():
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
    assert resp.status_code == 200

def test_register_email_conflict():
    resp = requests.post(
        url="http://localhost/auth/register",
        json={
            "email": "example@example.com",
            "password": "strong_pass_123",
        },
    )
    
    requests.post(
        url="http://localhost/auth/email/confirm",
        json={
            "user_id": resp.json()["data"]["id"],
            "code": "0000000000",
        },
    )

    resp = requests.post(
        url="http://localhost/auth/register",
        json={
            "email": "example@example.com",
            "password": "strong_pass_123",
        },
    )
    assert resp.status_code == 409

def test_register_confirm_email_with_wrong_code():
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
            "code": "0000000001",
        },
    )
    
    assert resp.status_code == 400

def test_register_confirm_email_when_attempts_exceeded():
    resp = requests.post(
        url="http://localhost/auth/register",
        json={
            "email": "example@example.com",
            "password": "strong_pass_123",
        },
    )
    
    user_id = resp.json()["data"]["id"]
    
    for _ in range(3):
        resp = requests.post(
            url="http://localhost/auth/email/confirm",
            json={
                "user_id": user_id,
                "code": "0000000001",
            },
        )
        assert resp.status_code == 400
    
    resp = requests.post(
        url="http://localhost/auth/email/confirm",
        json={
            "user_id": user_id,
            "code": "0000000001",
        },
    )
    assert resp.status_code == 403

def test_register_resend_code():
    resp = requests.post(
        url="http://localhost/auth/register",
        json={
            "email": "example@example.com",
            "password": "strong_pass_123",
        },
    )
    
    user_id = resp.json()["data"]["id"]
    
    for _ in range(3):
        resp = requests.post(
            url="http://localhost/auth/email/confirm",
            json={
                "user_id": user_id,
                "code": "0000000001",
            },
        )
        assert resp.status_code == 400
    
    resp = requests.post(
        url="http://localhost/auth/email/resend",
        json={
            "user_id": user_id,
        },
    )
    assert resp.status_code == 200
    
    requests.post(
        url="http://localhost/auth/email/confirm",
        json={
            "user_id": user_id,
            "code": "0000000000",
        },
    )
    assert resp.status_code == 200
